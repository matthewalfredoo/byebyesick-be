package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/appvalidator"
	"halodeksik-be/app/dto"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/dto/requestdto"
	"halodeksik-be/app/dto/responsedto"
	"halodeksik-be/app/dto/uriparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/usecase"
	"halodeksik-be/app/ws"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	HandshakeTimeout: appconstant.DefaultRequestTimeout * time.Second,
}

type ChatHandler struct {
	hub                   *ws.Hub
	consultationSessionUC usecase.ConsultationSessionUseCase
	consultationMessageUC usecase.ConsultationMessageUseCase
	profileUC             usecase.ProfileUseCase
	validator             appvalidator.AppValidator
}

func NewChatHandler(
	hub *ws.Hub,
	consultationSessionUC usecase.ConsultationSessionUseCase,
	consultationMessageUC usecase.ConsultationMessageUseCase,
	profileUC usecase.ProfileUseCase,
	validator appvalidator.AppValidator,
) *ChatHandler {
	return &ChatHandler{
		hub:                   hub,
		consultationSessionUC: consultationSessionUC,
		consultationMessageUC: consultationMessageUC,
		profileUC:             profileUC,
		validator:             validator,
	}
}

func (h *ChatHandler) CreateRoom(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	req := requestdto.AddConsultationSession{}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	err = h.validator.Validate(req)
	if err != nil {
		return
	}

	addedOrFound, err := h.consultationSessionUC.Add(ctx, req.ToConsultationSessionUseCase())
	if err != nil && errors.Is(err, apperror.ErrChatStillOngoing) {
		roomId := addedOrFound.Id
		_, isRoomExisted := h.hub.ConsultationSessions[roomId]
		if !isRoomExisted {
			h.hub.ConsultationSessions[roomId] = &ws.ConsultationSession{
				Id:        roomId,
				DoctorId:  addedOrFound.DoctorId,
				PatientId: addedOrFound.UserId,
				Clients:   make(map[int64]*ws.Client),
			}
		}
		return
	}

	if err != nil && !errors.Is(err, apperror.ErrChatStillOngoing) {
		return
	}

	roomId := addedOrFound.Id
	h.hub.ConsultationSessions[roomId] = &ws.ConsultationSession{
		Id:        roomId,
		DoctorId:  addedOrFound.DoctorId,
		PatientId: addedOrFound.UserId,
		Clients:   make(map[int64]*ws.Client),
	}

	resp := dto.ResponseDto{Data: addedOrFound.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ChatHandler) GetById(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	sessionId := uri.Id
	sessionDb, err := h.consultationSessionUC.GetById(ctx, sessionId)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: sessionDb.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ChatHandler) JoinRoom(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	sessionId := uri.Id
	sessionDb, err := h.consultationSessionUC.GetById(ctx, sessionId)
	if err != nil {
		return
	}

	if sessionDb.ConsultationSessionStatusId != appconstant.ConsultationSessionStatusOngoing {
		err = apperror.ErrChatAlreadyEnded
		return
	}

	_, isRoomExisted := h.hub.ConsultationSessions[sessionId]
	if !isRoomExisted {
		h.hub.ConsultationSessions[sessionId] = &ws.ConsultationSession{
			Id:        sessionId,
			DoctorId:  sessionDb.DoctorId,
			PatientId: sessionDb.UserId,
			Clients:   make(map[int64]*ws.Client),
		}
	}

	clientIdCtx := ctx.Request.Context().Value(appconstant.ContextKeyUserId)
	clientId := clientIdCtx.(int64)

	roleIdCtx := ctx.Request.Context().Value(appconstant.ContextKeyRoleId)
	roleId := roleIdCtx.(int64)

	var user *entity.User

	if roleId == appconstant.UserRoleIdDoctor {
		user, err = h.profileUC.GetDoctorProfileByUserId(ctx, sessionDb.DoctorId)
		if err != nil {
			return
		}
	}

	if roleId == appconstant.UserRoleIdUser {
		user, err = h.profileUC.GetUserProfileByUserId(ctx, sessionDb.UserId)
		if err != nil {
			return
		}
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		Conn:      conn,
		Message:   make(chan *responsedto.WsConsultationMessage, 10),
		SenderId:  clientId,
		SessionId: sessionId,
		Profile:   user.GetProfile(),
	}

	h.hub.Register <- client

	go client.WriteMessage()
	go client.ReadMessage(h.hub, h.consultationMessageUC, h.consultationSessionUC)
}

func (h *ChatHandler) GetAllByUserIdOrDoctorId(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	query := queryparamdto.GetAllConsultationSessions{}
	err = ctx.ShouldBindQuery(&query)
	if err != nil {
		return
	}

	err = h.validator.Validate(query)
	if err != nil {
		return
	}

	param := query.ToGetAllParams()
	paginatedItems, err := h.consultationSessionUC.GetAllByUserIdOrDoctorId(ctx, param)
	if err != nil {
		return
	}

	resps := make([]*responsedto.ConsultationSessionResponse, 0)
	for _, session := range paginatedItems.Items.([]*entity.ConsultationSession) {
		resps = append(resps, session.ToResponse())
	}
	paginatedItems.Items = resps

	resp := dto.ResponseDto{Data: paginatedItems}
	ctx.JSON(http.StatusOK, resp)
}

func (h *ChatHandler) EditStatusAsEnded(ctx *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			err = WrapError(err)
			_ = ctx.Error(err)
		}
	}()

	uri := uriparamdto.ResourceById{}
	err = ctx.ShouldBindUri(&uri)
	if err != nil {
		return
	}

	sessionId := uri.Id
	edited, err := h.consultationSessionUC.EditStatusAsEnded(ctx, sessionId)
	if err != nil {
		return
	}

	resp := dto.ResponseDto{Data: edited.ToResponse()}
	ctx.JSON(http.StatusOK, resp)
}

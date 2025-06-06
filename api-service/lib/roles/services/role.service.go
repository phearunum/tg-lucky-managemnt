// Role_service.go
package roles

import (
	dto "api-service/lib/roles/dto"
	"api-service/lib/roles/repository"
	"api-service/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) HandleRoleListRequest(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {

	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}

	roles, total, err := s.repo.GetRoleList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {
		utils.ErrorLog(nil, err.Error())

	}

	var roleDTOs []*dto.RoleDTO
	for _, role := range roles {
		roleDTOs = append(roleDTOs, &dto.RoleDTO{
			ID:         int(role.ID),
			RoleName:   role.RoleName,
			RoleStatus: role.RoleStatus,
			RoleKey:    role.RoleKey,
		})
	}
	return utils.NewServicePaginationResponse(roleDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "")

}
func (s *RoleService) HandleRoleListRequest1(obj []byte) {
	var requestDto utils.PaginationRequestDTO
	err := json.Unmarshal(obj, &requestDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		//return nil, 0, fmt.Errorf("error retrieving user list: %v", err)
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}
	roles, total, err := s.repo.GetRoleList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {

		utils.ErrorLog(nil, err.Error())
	}
	// Map roles (models) to DTOs
	var roleDTOs []*dto.RoleDTO
	for _, role := range roles {
		roleDTOs = append(roleDTOs, &dto.RoleDTO{
			ID:         int(role.ID),
			RoleName:   role.RoleName,
			RoleStatus: role.RoleStatus,
			RoleKey:    role.RoleKey,
		})
	}
	utils.NewServicePaginationResponse(roleDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, "Success", logrus.InfoLevel, "RoleService [HandleRoleListRequest]")

}
func (us *RoleService) HandleRoleCreateRequest(createRoleRequest dto.RoleCreateDTO) (*dto.RoleDTO, error) {
	createdRole, err := us.repo.CreateRole(&createRoleRequest)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	roleDTO := &dto.RoleDTO{
		RoleName:   createdRole.RoleName,
		RoleStatus: createdRole.RoleStatus,
		RoleKey:    createdRole.RoleKey,
	}
	utils.InfoLog(roleDTO, string(utils.SuccessMessage))
	return roleDTO, nil

}

func (us *RoleService) HandleRoleUpdateRequest(updateRoleRequest dto.RoleUpdateDTO) (*dto.RoleDTO, error) {

	RoleID := uint(updateRoleRequest.ID)
	updatedRole, err := us.repo.UpdateRole(RoleID, &updateRoleRequest)
	if err != nil {
		utils.ErrorLog(nil, "Failed to update role: "+err.Error())
		return nil, err
	}

	roleDTO := &dto.RoleDTO{
		ID:         int(updatedRole.ID),
		RoleName:   updatedRole.RoleName,
		RoleStatus: updatedRole.RoleStatus,
		RoleKey:    updatedRole.RoleKey,
	}

	utils.InfoLog(roleDTO, string(utils.SuccessMessage))
	return roleDTO, nil
}

func (us *RoleService) HandleRoleByIdRequest(RoleID int) (*dto.RoleDTO, error) {

	RoleFiller, err := us.repo.GetRoleByID(uint(RoleID))
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	roleDTO := &dto.RoleDTO{
		ID:         int(RoleFiller.ID),
		RoleName:   RoleFiller.RoleName,
		RoleStatus: RoleFiller.RoleStatus,
		RoleKey:    RoleFiller.RoleKey,
	}

	utils.InfoLog(roleDTO, string(utils.SuccessMessage))
	return roleDTO, nil

}

func (us *RoleService) HandleRoleDeleteRequest(obj []byte) (*dto.RoleDTO, error) {
	var deleteRoleRequest dto.RoleDTO
	if err := json.Unmarshal(obj, &deleteRoleRequest); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	err := us.repo.DeleteRoleByID(uint(deleteRoleRequest.ID))
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	utils.InfoLog(nil, string(utils.SuccessMessage))
	return nil, nil
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kifangamukundi/gm/loan/controllers"
	"github.com/kifangamukundi/gm/loan/loanrepository"
	"github.com/kifangamukundi/gm/loan/models"
	"github.com/kifangamukundi/gm/loan/services"
	"gorm.io/gorm"
)

func InitializeRoutes(r *gin.Engine, db *gorm.DB) {
	// repo layer
	loanRepo := loanrepository.NewLoanRepository(db)

	// service layer
	service := services.NewEntityService(loanRepo)

	// Models layer
	userModel := models.NewUserModel(service)
	roleModel := models.NewRoleModel(service)
	permissionModel := models.NewPermissionModel(service)
	countryModel := models.NewCountryModel(service)
	regionModel := models.NewRegionModel(service)

	countyModel := models.NewCountyModel(service)
	subCountyModel := models.NewSubCountyModel(service)
	wardModel := models.NewWardModel(service)
	locationModel := models.NewLocationModel(service)
	subLocationModel := models.NewSubLocationModel(service)
	villageModel := models.NewVillageModel(service)
	roadModel := models.NewRoadModel(service)
	plotModel := models.NewPlotModel(service)
	unitModel := models.NewUnitModel(service)

	agentModel := models.NewAgentModel(service)
	groupModel := models.NewGroupModel(service)
	officerModel := models.NewOfficerModel(service)
	memberModel := models.NewMemberModel(service)
	loanModel := models.NewLoanModel(service)
	disburseModel := models.NewDisburseModel(service)

	// Controllers layer
	userController := controllers.NewUserController(userModel)
	roleController := controllers.NewRoleController(roleModel)
	permissionController := controllers.NewPermissionController(permissionModel)
	countryController := controllers.NewCountryController(countryModel)
	regionController := controllers.NewRegionController(regionModel)
	countyController := controllers.NewCountyController(countyModel)
	subCountyController := controllers.NewSubCountyController(subCountyModel)
	wardController := controllers.NewWardController(wardModel)
	locationController := controllers.NewLocationController(locationModel)
	subLocationController := controllers.NewSubLocationController(subLocationModel)
	villageController := controllers.NewVillageController(villageModel)
	roadController := controllers.NewRoadController(roadModel)
	plotController := controllers.NewPlotController(plotModel)
	unitController := controllers.NewUnitController(unitModel)
	agentController := controllers.NewAgentController(agentModel, userModel)
	groupController := controllers.NewGroupController(groupModel, userModel, agentModel)
	officerController := controllers.NewOfficerController(officerModel, userModel)
	memberController := controllers.NewMemberController(memberModel, userModel, groupModel)
	loanController := controllers.NewLoanController(loanModel, disburseModel, userModel, officerModel, agentModel, groupModel, memberModel)

	UserRoutes(r, userController, db)
	RoleRoutes(r, roleController, db)
	PermissionRoutes(r, permissionController, db)
	CountryRoutes(r, countryController, db)
	RegionRoutes(r, regionController, db)
	CountyRoutes(r, countyController, db)
	SubCountyRoutes(r, subCountyController, db)
	WardRoutes(r, wardController, db)
	LocationRoutes(r, locationController, db)
	SubLocationRoutes(r, subLocationController, db)
	VillageRoutes(r, villageController, db)
	RoadRoutes(r, roadController, db)
	PlotRoutes(r, plotController, db)
	UnitRoutes(r, unitController, db)
	AgentRoutes(r, agentController, db)
	GroupRoutes(r, groupController, db)
	OfficerRoutes(r, officerController, db)
	MemberRoutes(r, memberController, db)
	LoanRoutes(r, loanController, db)

	MediaRoutes(r, db)
}

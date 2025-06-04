package mongodb

import (
	"go.fork.vn/config"
	"go.fork.vn/di"
)

// ServiceProvider định nghĩa interface cho MongoDB service provider.
//
// ServiceProvider kế thừa từ di.ServiceProvider và định nghĩa
// các phương thức cần thiết cho một MongoDB service provider.
type ServiceProvider interface {
	di.ServiceProvider
}

// serviceProvider là implementation của ServiceProvider.
//
// serviceProvider chịu trách nhiệm đăng ký các dịch vụ MongoDB vào DI container
// và cung cấp MongoDB client cho các module khác trong ứng dụng.
type serviceProvider struct {
	providers []string
}

// NewServiceProvider tạo một MongoDB service provider mới.
func NewServiceProvider() ServiceProvider {
	return &serviceProvider{}
}

// Register đăng ký các dịch vụ MongoDB với DI container.
//
// Phương thức này đăng ký MongoDB manager vào container DI của ứng dụng.
// Nó khởi tạo một MongoDB manager mới và đăng ký nó dưới các alias "mongodb",
// "mongo.client" và "mongo".
//
// Params:
//   - app: Interface của ứng dụng, phải cung cấp phương thức Container() để lấy container DI
func (p *serviceProvider) Register(app di.Application) {
	if app == nil {
		panic("application cannot be nil")
	}
	// Lấy container từ app
	c := app.Container()
	// Kiểm tra xem container có tồn tại không
	if c == nil {
		panic("container cannot be nil")
	}

	// Kiểm tra xem container đã có config manager chưa
	mongoConfig := DefaultConfig()
	configService, ok := c.MustMake("config").(config.Manager)
	if !ok {
		panic("MongoDB provider requires config service to be registered")
	}
	err := configService.UnmarshalKey("mongodb", &mongoConfig)
	if err != nil {
		panic("MongoDB config unmarshal error: " + err.Error())
	}

	// Tạo MongoDB manager với cấu hình
	manager := NewManagerWithConfig(*mongoConfig) // Đăng ký manager vào container
	c.Instance("mongodb", manager)

	// Đăng ký client và database instances (sẽ sử dụng lazy initialization trong manager)
	client := manager.Client()
	c.Instance("mongodb.client", client)

	database := manager.Database()
	c.Instance("mongodb.database", database)

	// Add to providers list - the test expects these specific entries
	p.providers = append(p.providers, "mongodb")
	p.providers = append(p.providers, "mongodb.client")
	p.providers = append(p.providers, "mongodb.database")

}

// Boot khởi động MongoDB provider.
//
// Phương thức này khởi động MongoDB provider sau khi tất cả các service provider đã được đăng ký.
// Trong trường hợp này, không cần thực hiện thêm tác vụ nào trong Boot vì các cấu hình
// đã được xử lý trong Register.
//
// Params:
//   - app: di.Application của ứng dụng
func (p *serviceProvider) Boot(app di.Application) {
	// Không cần thực hiện thêm tác vụ nào trong Boot
	// vì cấu hình đã được xử lý trong Register
	if app == nil {
		panic("application cannot be nil")
	}
}

// Providers trả về danh sách các service được cung cấp bởi MongoDB provider.
//
// Phương thức này trả về danh sách các abstract type mà MongoDB provider đăng ký với container.
// Danh sách này được sử dụng để kiểm tra dependencies và đảm bảo đúng thứ tự khởi tạo.
//
// Trả về:
//   - []string: danh sách các service được cung cấp
func (p *serviceProvider) Providers() []string {
	return p.providers
}

// Requires trả về danh sách các dependency mà MongoDB provider phụ thuộc.
//
// Trả về:
//   - []string: danh sách các service provider khác mà provider này yêu cầu
func (p *serviceProvider) Requires() []string {
	return []string{
		// MongoDB provider yêu cầu config provider để đọc cấu hình
		"config",
	}
}

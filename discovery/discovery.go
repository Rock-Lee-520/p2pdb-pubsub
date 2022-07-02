package discovery

type Discovery interface {
	// create a new libp2p Host that listens on a random TCP port
	Connect() (Host string) 
}

type DiscoveryFactory struct {
}

func NewDiscoveryFactory() *DiscoveryFactory {
	return &DiscoveryFactory{}
}

func (d *DiscoveryFactory) Connect() (Host string)  (Discovery, error){
	return d,nil
}

func (d *DiscoveryFactory) Create(host string) (Discovery, error) {
	return d,nil
}


func (d *DiscoveryFactory) Register(host string) (Discovery, error) {
	return d,nil
}

func (d *DiscoveryFactory) Disconnect() (Discovery, error) {
	return d,nil
}	



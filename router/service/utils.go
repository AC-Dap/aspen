package service

func (s Status) String() string {
	switch s {
	case NotInitialized:
		return "Not initialized"
	case Building:
		return "Building"
	case Built:
		return "Built"
	case Starting:
		return "Starting"
	case Started:
		return "Started"
	case Stopping:
		return "Stopping"
	case Stopped:
		return "Stopped"
	default:
		return "Unknown status"
	}
}

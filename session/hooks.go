package session

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

type BeforeInsertHook interface {
	BeforeInsert(*Session)
}
type AfterInsertHook interface {
	AfterInsert(*Session)
}
type BeforeUpdateHook interface {
	BeforeUpdate(*Session)
}
type AfterUpdateHook interface {
	AfterUpdate(*Session)
}
type BeforeDeleteHook interface {
	BeforeDelete(*Session)
}
type AfterDeleteHook interface {
	AfterDelete(*Session)
}
type BeforeQueryHook interface {
	BeforeQuery(*Session)
}
type AfterQueryHook interface {
	AfterQuery(*Session)
}

func (s *Session) CallMethod(method string, value interface{}) {
	switch method {
	case BeforeInsert:
		if hook, ok := value.(BeforeInsertHook); ok {
			hook.BeforeInsert(s)
		}
	case AfterInsert:
		if hook, ok := value.(AfterInsertHook); ok {
			hook.AfterInsert(s)
		}
	case BeforeUpdate:
		if hook, ok := value.(BeforeUpdateHook); ok {
			hook.BeforeUpdate(s)
		}
	case AfterUpdate:
		if hook, ok := value.(AfterUpdateHook); ok {
			hook.AfterUpdate(s)
		}
	case BeforeDelete:
		if hook, ok := value.(BeforeDeleteHook); ok {
			hook.BeforeDelete(s)
		}
	case AfterDelete:
		if hook, ok := value.(AfterDeleteHook); ok {
			hook.AfterDelete(s)
		}
	case BeforeQuery:
		if hook, ok := value.(BeforeQueryHook); ok {
			hook.BeforeQuery(s)
		}
	case AfterQuery:
		if hook, ok := value.(AfterQueryHook); ok {
			hook.AfterQuery(s)
		}
	}
	return
}

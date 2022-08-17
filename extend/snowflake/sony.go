package snowflake

import (
	"errors"

	"github.com/sony/sonyflake"
)

// sony 实现
type sony struct {
	SonyFlake *sonyflake.Sonyflake
}

// NewID NewID
func (s *sony) NewID() (uint64, error) {
	return s.SonyFlake.NextID()
}

// newSony 实例
func newSony() (*sony, error) {
	var st sonyflake.Settings
	//st.MachineID = awsutil.AmazonEC2MachineID
	sf := sonyflake.NewSonyflake(st)

	if sf == nil {
		return nil, errors.New("创建SonyFlake对象失败")
	}

	impl := &sony{}
	impl.SonyFlake = sf

	return impl, nil
}

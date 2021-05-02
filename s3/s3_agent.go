package s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// s3Agent is struct that agent API about aws s3 including put object, delete object, etc ...
type s3Agent struct {
	session *session.Session
}

func New(ses *session.Session) *s3Agent {
	return &s3Agent{
		session: ses,
	}
}

// PutObject method put(insert or update) object to s3
func (sa *s3Agent) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return s3.New(sa.session).PutObject(input)
}

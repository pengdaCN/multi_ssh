package m_terminal

import "golang.org/x/crypto/ssh"

type Result struct {
	code    int
	msg     string
	total   string
	stdout  string
	stderr  string
	errInfo string
	err     error
}

func (r *Result) Code() int {
	return r.code
}

func (r *Result) Msg() string {
	return r.msg
}

func (r *Result) TotalMsg() string {
	return r.total
}

func (r *Result) Stdout() string {
	return r.stdout
}

func (r *Result) Stderr() string {
	return r.stderr
}

func (r *Result) ErrInfo() string {
	return r.errInfo
}

func (r *Result) Err() error {
	return r.err
}

func buildErrWithText(txt string) *Result {
	r := new(Result)
	r.msg = txt
	r.errInfo = txt
	return r
}

func buildRst(msg string, err error) *Result {
	r := new(Result)
	r.msg = msg
	r.err = err
	if err != nil {
		r.errInfo = err.Error()
	}
	if exit, ok := err.(*ssh.ExitError); ok {
		r.code = exit.ExitStatus()
	}
	return r
}

func buildRstByErr(err error) *Result {
	rst := buildRst("", err)
	if err == nil {
		rst.msg = "OK"
	}
	return rst
}

func buildRstWithOK() *Result {
	return buildRstByErr(nil)
}

func BuildRstWithCode(code int) *Result {
	r := new(Result)
	r.code = code
	if code != 0 {
		r.msg = "ERROR"
		r.errInfo = r.msg
	} else {
		r.msg = "OK"
	}
	return r
}

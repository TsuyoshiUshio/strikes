package ui

type Process interface {
	PrintQuestion() error
	WaitForInput() (string, error)
	Validate(answer string) bool
	IsTargetParameterFilled(parameter interface{}) bool
	UpdateParameter(answer string, parameter interface{}) (interface{}, error)
	ShowValidateError(answer string)
	SetNext(process *Process)
	SetParameter(parameter interface{})
	Next() *Process
}

func Execute(p *Process, parameter interface{}) (interface{}, error) {
	(*p).SetParameter(parameter)
	if (*p).IsTargetParameterFilled(parameter) {
		if (*p).Next() != nil {
			return Execute((*p).Next(), parameter)
		} else {
			return parameter, nil
		}
	}
	err := (*p).PrintQuestion()
	if err != nil {
		return nil, err
	}
	answer, err := (*p).WaitForInput()
	if err != nil {
		return nil, err
	}

	if (*p).Validate(answer) {
		result, err := (*p).UpdateParameter(answer, parameter)
		if err != nil {
			return nil, err
		}
		if (*p).Next() != nil {
			return Execute((*p).Next(), parameter)
		} else {
			return result, nil
		}
	} else {
		(*p).ShowValidateError(answer)
		return Execute(p, parameter)
	}
}

type ProcessBuilder struct {
	Head    *Process
	Current *Process
}

func NewProcessBuilder() *ProcessBuilder {
	return &ProcessBuilder{}
}
func (b *ProcessBuilder) Append(process *Process) {
	if b.Head == nil {
		b.Head = process
	}
	if b.Current == nil {
		b.Current = process
	} else {
		(*b.Current).SetNext(process)
		b.Current = process
	}
}

func (b *ProcessBuilder) Build() *Process {
	return b.Head
}

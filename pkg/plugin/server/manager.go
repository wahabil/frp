// Copyright 2019 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plugin

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fatedier/frp/pkg/util/util"
	"github.com/fatedier/frp/pkg/util/xlog"
)

type Manager struct {
	loginPlugins       []Plugin
	exitPlugins        []Plugin
	newProxyPlugins    []Plugin
	closeProxyPlugins  []Plugin
	pingPlugins        []Plugin
	newWorkConnPlugins []Plugin
	newUserConnPlugins []Plugin
}

func NewManager() *Manager {
	return &Manager{
		loginPlugins:       make([]Plugin, 0),
		exitPlugins:        make([]Plugin, 0),
		newProxyPlugins:    make([]Plugin, 0),
		closeProxyPlugins:  make([]Plugin, 0),
		pingPlugins:        make([]Plugin, 0),
		newWorkConnPlugins: make([]Plugin, 0),
		newUserConnPlugins: make([]Plugin, 0),
	}
}

func (m *Manager) Register(p Plugin) {
	if p.IsSupport(OpLogin) {
		m.loginPlugins = append(m.loginPlugins, p)
	}
	if p.IsSupport(OpExit) {
		m.exitPlugins = append(m.exitPlugins, p)
	}
	if p.IsSupport(OpNewProxy) {
		m.newProxyPlugins = append(m.newProxyPlugins, p)
	}
	if p.IsSupport(OpCloseProxy) {
		m.closeProxyPlugins = append(m.closeProxyPlugins, p)
	}
	if p.IsSupport(OpPing) {
		m.pingPlugins = append(m.pingPlugins, p)
	}
	if p.IsSupport(OpNewWorkConn) {
		m.newWorkConnPlugins = append(m.newWorkConnPlugins, p)
	}
	if p.IsSupport(OpNewUserConn) {
		m.newUserConnPlugins = append(m.newUserConnPlugins, p)
	}
}

func (m *Manager) Login(content *LoginContent) (*LoginContent, error) {
	if len(m.loginPlugins) == 0 {
		return content, nil
	}

	var (
		res = &Response{
			Reject:   false,
			Unchange: true,
		}
		retContent interface{}
		err        error
	)
	reqid, _ := util.RandID()
	xl := xlog.New().AppendPrefix("reqid: " + reqid)
	ctx := xlog.NewContext(context.Background(), xl)
	ctx = NewReqidContext(ctx, reqid)

	for _, p := range m.loginPlugins {
		res, retContent, err = p.Handle(ctx, OpLogin, *content)
		if err != nil {
			xl.Warn("send Login request to plugin [%s] error: %v", p.Name(), err)
			return nil, errors.New("send Login request to plugin error")
		}
		if res.Reject {
			return nil, fmt.Errorf("%s", res.RejectReason)
		}
		if !res.Unchange {
			content = retContent.(*LoginContent)
		}
	}
	return content, nil
}

func (m *Manager) Exit(content *ExitContent) error {
	if len(m.exitPlugins) == 0 {
		return nil
	}

	errs := make([]string, 0)
	reqid, _ := util.RandID()
	xl := xlog.New().AppendPrefix("reqid: " + reqid)
	ctx := xlog.NewContext(context.Background(), xl)
	ctx = NewReqidContext(ctx, reqid)

	for _, p := range m.exitPlugins {
		_, _, err := p.Handle(ctx, OpExit, *content)
		if err != nil {
			xl.Warn("send Exit request to plugin [%s] error: %v", p.Name(), err)
			errs = append(errs, fmt.Sprintf("[%s]: %v", p.Name(), err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("send Exit request to plugin errors: %s", strings.Join(errs, "; "))
	} else {
		return nil
	}
}

func (m *Manager) NewProxy(content *NewProxyContent) (*NewProxyContent, error) {
	if len(m.newProxyPlugins) == 0 {
		return content, nil
	}

	var (
		res = &Response{
			Reject:   false,
			Unchange: true,
		}
		retContent interface{}
		err        error
	)
	reqid, _ := util.RandID()
	xl := xlog.New().AppendPrefix("reqid: " + reqid)
	ctx := xlog.NewContext(context.Background(), xl)
	ctx = NewReqidContext(ctx, reqid)

	for _, p := range m.newProxyPlugins {
		res, retContent, err = p.Handle(ctx, OpNewProxy, *content)
		if err != nil {
			xl.Warn("send NewProxy request to plugin [%s] error: %v", p.Name(), err)
			return nil, errors.New("send NewProxy request to plugin error")
		}
		if res.Reject {
			return nil, fmt.Errorf("%s", res.RejectReason)
		}
		if !res.Unchange {
			content = retContent.(*NewProxyContent)
		}
	}
	return content, nil
}

func (m *Manager) CloseProxy(content *CloseProxyContent) error {
	if len(m.closeProxyPlugins) == 0 {
		return nil
	}

	errs := make([]string, 0)
	reqid, _ := util.RandID()
	xl := xlog.New().AppendPrefix("reqid: " + reqid)
	ctx := xlog.NewContext(context.Background(), xl)
	ctx = NewReqidContext(ctx, reqid)

	for _, p := range m.closeProxyPlugins {
		_, _, err := p.Handle(ctx, OpCloseProxy, *content)
		if err != nil {
			xl.Warn("send CloseProxy request to plugin [%s] error: %v", p.Name(), err)
			errs = append(errs, fmt.Sprintf("[%s]: %v", p.Name(), err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("send CloseProxy request to plugin errors: %s", strings.Join(errs, "; "))
	} else {
		return nil
	}
}

func (m *Manager) Ping(content *PingContent) (*PingContent, error) {
	if len(m.pingPlugins) == 0 {
		return content, nil
	}

	var (
		res = &Response{
			Reject:   false,
			Unchange: true,
		}
		retContent interface{}
		err        error
	)
	reqid, _ := util.RandID()
	xl := xlog.New().AppendPrefix("reqid: " + reqid)
	ctx := xlog.NewContext(context.Background(), xl)
	ctx = NewReqidContext(ctx, reqid)

	for _, p := range m.pingPlugins {
		res, retContent, err = p.Handle(ctx, OpPing, *content)
		if err != nil {
			xl.Warn("send Ping request to plugin [%s] error: %v", p.Name(), err)
			return nil, errors.New("send Ping request to plugin error")
		}
		if res.Reject {
			return nil, fmt.Errorf("%s", res.RejectReason)
		}
		if !res.Unchange {
			content = retContent.(*PingContent)
		}
	}
	return content, nil
}

func (m *Manager) NewWorkConn(content *NewWorkConnContent) (*NewWorkConnContent, error) {
	if len(m.newWorkConnPlugins) == 0 {
		return content, nil
	}

	var (
		res = &Response{
			Reject:   false,
			Unchange: true,
		}
		retContent interface{}
		err        error
	)
	reqid, _ := util.RandID()
	xl := xlog.New().AppendPrefix("reqid: " + reqid)
	ctx := xlog.NewContext(context.Background(), xl)
	ctx = NewReqidContext(ctx, reqid)

	for _, p := range m.newWorkConnPlugins {
		res, retContent, err = p.Handle(ctx, OpPing, *content)
		if err != nil {
			xl.Warn("send NewWorkConn request to plugin [%s] error: %v", p.Name(), err)
			return nil, errors.New("send NewWorkConn request to plugin error")
		}
		if res.Reject {
			return nil, fmt.Errorf("%s", res.RejectReason)
		}
		if !res.Unchange {
			content = retContent.(*NewWorkConnContent)
		}
	}
	return content, nil
}

func (m *Manager) NewUserConn(content *NewUserConnContent) (*NewUserConnContent, error) {
	if len(m.newUserConnPlugins) == 0 {
		return content, nil
	}

	var (
		res = &Response{
			Reject:   false,
			Unchange: true,
		}
		retContent interface{}
		err        error
	)
	reqid, _ := util.RandID()
	xl := xlog.New().AppendPrefix("reqid: " + reqid)
	ctx := xlog.NewContext(context.Background(), xl)
	ctx = NewReqidContext(ctx, reqid)

	for _, p := range m.newUserConnPlugins {
		res, retContent, err = p.Handle(ctx, OpNewUserConn, *content)
		if err != nil {
			xl.Info("send NewUserConn request to plugin [%s] error: %v", p.Name(), err)
			return nil, errors.New("send NewUserConn request to plugin error")
		}
		if res.Reject {
			return nil, fmt.Errorf("%s", res.RejectReason)
		}
		if !res.Unchange {
			content = retContent.(*NewUserConnContent)
		}
	}
	return content, nil
}

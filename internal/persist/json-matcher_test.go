package persist_test

import (
	"fmt"

	"github.com/onsi/gomega/types"
	"github.com/snivilised/traverse/internal/opts/json"
	"github.com/snivilised/traverse/internal/persist"

	"github.com/snivilised/traverse/pref"
)

type MarshalJSONMatcher struct {
	o   *pref.Options
	err error
}

func HaveMarshaledEqual(o *pref.Options) types.GomegaMatcher {
	return &MarshalJSONMatcher{
		o: o,
	}
}

func (m *MarshalJSONMatcher) Match(actual interface{}) (bool, error) {
	jo, ok := actual.(*json.Options)

	if !ok {
		return false, fmt.Errorf("❌ matcher expected a *json.Options instance (%T)", jo)
	}

	if err := (&persist.Comparison{
		O:  m.o,
		JO: jo,
	}).Equals(); err != nil {
		m.err = err
		return false, err
	}

	return true, nil
}

func (m *MarshalJSONMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected\n\t%v\nJSON Marshal conversion result in equal result", m.err)
}

func (m *MarshalJSONMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected\n\t%v\nJSON Marshal conversion result in NON equal result", m.err)
}

type UnMarshalJSONMatcher struct {
	jo  *json.Options
	err error
}

func HaveUnMarshaledEqual(jo *json.Options) types.GomegaMatcher {
	return &UnMarshalJSONMatcher{
		jo: jo,
	}
}

func (m *UnMarshalJSONMatcher) Match(actual interface{}) (bool, error) {
	o, ok := actual.(*pref.Options)

	if !ok {
		return false, fmt.Errorf("❌ matcher expected a *pref.Options instance (%T)", o)
	}

	if err := (&persist.Comparison{
		O:  o,
		JO: m.jo,
	}).Equals(); err != nil {
		m.err = err
		return false, err
	}

	return true, nil
}

func (m *UnMarshalJSONMatcher) FailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected\n\t%v\nJSON UnMarshal conversion result in equal result", m.err)
}

func (m *UnMarshalJSONMatcher) NegatedFailureMessage(_ interface{}) string {
	return fmt.Sprintf("🔥 Expected\n\t%v\nJSON UnMarshal conversion result in NON equal result", m.err)
}

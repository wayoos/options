package opts

type RuleSet struct {
	rules map[string]*rule
}

type rule struct {
	optionID  string
	deps      []*rule
	conflicts []*rule
}

func (r *rule) addDep(to *rule) {
	r.deps = append(r.deps, to)
}

func (r *rule) addConflict(to *rule) {
	r.conflicts = append(r.conflicts, to)
}

func (c *RuleSet) AddDep(from, to string) {
	fromRule := c.findOrCreate(from)
	toRule := c.findOrCreate(to)
	fromRule.addDep(toRule)
}

func (c *RuleSet) AddConflict(from, to string) {
	fromRule := c.findOrCreate(from)
	toRule := c.findOrCreate(to)
	fromRule.addConflict(toRule)
}

func (c *RuleSet) findOrCreate(optionID string) *rule {
	r, present := c.rules[optionID]
	if present {
		return r
	}
	newRule := rule{optionID: optionID}
	c.rules[optionID] = &newRule
	return &newRule
}

func (c *RuleSet) IsCoherent() bool {

	for _, fromRule := range c.rules {
		for _, toConflictRule := range fromRule.conflicts {

			for _, checkRule := range c.rules {
				_, hasFromDep := findDep(checkRule, fromRule.optionID, []string{})
				_, hasToDep := findDep(checkRule, toConflictRule.optionID, []string{})
				if hasFromDep && hasToDep {
					return false
				}
			}
		}
	}

	return true
}

func (c *RuleSet) findConflictsTo(id string) []string {
	conflictsList := []string{}

	for _, conflictTo := range c.rules[id].conflicts {
		conflictsList = append(conflictsList, conflictTo.optionID)
	}

	for _, fromRule := range c.rules {
		for _, toConflictRule := range fromRule.conflicts {
			if toConflictRule.optionID == id {
				conflictsList = append(conflictsList, fromRule.optionID)
			}
		}
	}
	return conflictsList
}

func findDep(r *rule, to string, detectLoop []string) (*rule, bool) {
	if r.optionID == to {
		return r, true
	}

	if len(r.deps) == 0 {
		return &rule{}, false
	}

	if contains(detectLoop, r.optionID) {
		// dependency loop
		return &rule{}, false
	}

	detectLoop = append(detectLoop, r.optionID)

	for _, depRule := range r.deps {
		res, present := findDep(depRule, to, detectLoop)
		if present {
			return res, true
		}
	}
	return &rule{}, false
}

func findDeps(r *rule, depList []string) []string {
	if contains(depList, r.optionID) {
		return depList
	}
	newDepList := append(depList, r.optionID)
	for _, depRule := range r.deps {
		newDepList = findDeps(depRule, newDepList)
	}
	return newDepList
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func merge(s1, s2 []string) []string {
	result := s1
	for _, a := range s2 {
		if !contains(result, a) {
			result = append(result, a)
		}
	}
	return result
}

// remove s2 element to s1
func remove(s1, s2 []string) []string {
	result := []string{}
	for _, a := range s1 {
		if !contains(s2, a) {
			result = append(result, a)
		}
	}
	return result
}

type Selection struct {
	ruleSet  *RuleSet
	selected []string
}

func (s *Selection) Toggle(id string) {
	rule := s.ruleSet.rules[id]

	// first find new toggled option
	newToggledIds := findDeps(rule, []string{})

	if contains(s.selected, id) {
		// unset

		s.selected = remove(s.selected, newToggledIds)
	} else {
		// set

		// remove conflicts
		for _, n := range newToggledIds {
			conflicts := s.ruleSet.findConflictsTo(n)
			for _, c := range conflicts {
				conflictRule := s.ruleSet.rules[c]
				conflictsIds := findDeps(conflictRule, []string{})
				s.selected = remove(s.selected, conflictsIds)
			}
		}

		s.selected = merge(newToggledIds, s.selected)
	}

}

func (c *Selection) StringSlice() []string {
	return c.selected
}

func NewRuleSet() *RuleSet {
	r := &RuleSet{}
	r.rules = make(map[string]*rule)
	return r
}

func New(r *RuleSet) *Selection {
	return &Selection{ruleSet: r}
}

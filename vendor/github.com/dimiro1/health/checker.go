package health

// Checker is a interface used to provide an indication of application health.
type Checker interface {
	Check() Health
}

// CheckerFunc is an adapter to allow the use of
// ordinary go functions as Checkers.
type CheckerFunc func() Health

func (f CheckerFunc) Check() Health {
	return f()
}

type checkerItem struct {
	name    string
	checker Checker
}

// CompositeChecker aggregate a list of Checkers
type CompositeChecker struct {
	checkers []checkerItem
	info     map[string]interface{}
}

// NewCompositeChecker creates a new CompositeChecker
func NewCompositeChecker() CompositeChecker {
	return CompositeChecker{}
}

// AddInfo adds a info value to the Info map
func (c *CompositeChecker) AddInfo(key string, value interface{}) *CompositeChecker {
	if c.info == nil {
		c.info = make(map[string]interface{})
	}

	c.info[key] = value

	return c
}

// AddChecker add a Checker to the aggregator
func (c *CompositeChecker) AddChecker(name string, checker Checker) {
	c.checkers = append(c.checkers, checkerItem{name: name, checker: checker})
}

// Check returns the combination of all checkers added
// if some check is not up, the combined is marked as down
func (c CompositeChecker) Check() Health {
	health := NewHealth()
	health.Up()

	healths := make(map[string]interface{})

	for _, item := range c.checkers {
		h := item.checker.Check()

		if !h.IsUp() && !health.IsDown() {
			health.Down()
		}

		healths[item.name] = h
	}

	health.info = healths

	// Extra Info
	for key, value := range c.info {
		health.AddInfo(key, value)
	}

	return health
}

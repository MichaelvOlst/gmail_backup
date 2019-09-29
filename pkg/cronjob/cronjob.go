package cronjob

// // Cronjob is a wrapper for robfig/cron
// type Cronjob struct {
// 	*cron.Cron
// }

// // New returs a new Cronjob
// func New() *Cronjob {
// 	c := cron.New()
// 	return &Cronjob{c}
// }

// // Add runs the Cronjob
// func (c *Cronjob) Add(expression string, fn func()) error {
// 	_, err := c.AddFunc(expression, fn)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Start runs the Cronjob
// func (c *Cronjob) Start() {
// 	c.Start()
// }

// // Stop stops the Cronjob
// func (c *Cronjob) Stop() {
// 	c.Stop()
// }

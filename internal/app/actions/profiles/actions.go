package profiles

import (
	"fmt"

	"github.com/alex-dev-master/etcdtui/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// RefreshProfileList reloads profiles into the list.
func (s *State) RefreshProfileList() {
	s.profileList.Clear()

	profiles := s.configManager.GetProfiles()

	if len(profiles) == 0 {
		s.profileList.AddItem("No profiles configured", "Press 'n' to create one", 0, nil)
		s.detailsView.SetText("[yellow]No profiles[-]\n\nCreate a new profile to get started.")
		return
	}

	for _, p := range profiles {
		profile := p // capture for closure
		secondaryText := p.Endpoints[0]
		if p.Default {
			secondaryText += " [default]"
		}

		s.profileList.AddItem(p.Name, secondaryText, 0, func() {
			s.Connect(profile)
		})
	}

	// Set up change handler to show details
	s.profileList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		if index < len(profiles) {
			s.ShowProfileDetails(profiles[index])
			s.selectedProfile = profiles[index]
		}
	})

	// Show first profile details
	if len(profiles) > 0 {
		s.ShowProfileDetails(profiles[0])
		s.selectedProfile = profiles[0]
	}
}

// ShowProfileDetails displays profile details in the details view.
func (s *State) ShowProfileDetails(p *config.Profile) {
	text := "[yellow::b]" + p.Name + "[-:-:-]\n\n"

	text += "[cyan]Endpoints:[-]\n"
	for _, ep := range p.Endpoints {
		text += "  • " + ep + "\n"
	}

	if p.HasAuth() {
		text += "\n[cyan]Authentication:[-]\n"
		text += "  Username: " + p.Username + "\n"
		if p.Password != "" {
			text += "  Password: ****\n"
		}
	}

	if p.HasTLS() {
		text += "\n[cyan]TLS:[-]\n"
		if p.TLS.CAFile != "" {
			text += "  CA: " + p.TLS.CAFile + "\n"
		}
		if p.TLS.CertFile != "" {
			text += "  Cert: " + p.TLS.CertFile + "\n"
		}
		if p.TLS.KeyFile != "" {
			text += "  Key: " + p.TLS.KeyFile + "\n"
		}
		if p.TLS.InsecureSkipVerify {
			text += "  [red]Insecure: skip verify[-]\n"
		}
	}

	if p.Default {
		text += "\n[green]✓ Default profile[-]"
	}

	s.detailsView.SetText(text)
}

// ShowNewProfileForm displays the form to create a new profile.
func (s *State) ShowNewProfileForm(restoreInput func()) {
	s.showProfileForm(nil, restoreInput)
}

// ShowEditProfileForm displays the form to edit an existing profile.
func (s *State) ShowEditProfileForm(p *config.Profile, restoreInput func()) {
	s.showProfileForm(p, restoreInput)
}

// closeForm closes the form and restores the main view.
func (s *State) closeForm(restoreInput func()) {
	s.app.SetRoot(s.rootFlex, true)
	if restoreInput != nil {
		restoreInput()
	}
}

// showProfileForm displays a form for creating/editing a profile.
func (s *State) showProfileForm(existing *config.Profile, restoreInput func()) {
	// Disable global input capture while form is active
	s.app.SetInputCapture(nil)

	form := tview.NewForm()

	// Default values
	name := ""
	endpoints := "localhost:2379"
	username := ""
	password := ""
	tlsEnabled := false
	caFile := ""
	certFile := ""
	keyFile := ""
	isDefault := false

	if existing != nil {
		name = existing.Name
		if len(existing.Endpoints) > 0 {
			endpoints = existing.Endpoints[0]
		}
		username = existing.Username
		password = existing.DecodePassword()
		isDefault = existing.Default
		if existing.TLS != nil {
			tlsEnabled = existing.TLS.Enabled
			caFile = existing.TLS.CAFile
			certFile = existing.TLS.CertFile
			keyFile = existing.TLS.KeyFile
		}
	}

	title := " New Profile "
	if existing != nil {
		title = " Edit Profile "
	}

	form.AddInputField("Name", name, 40, nil, nil)
	form.AddInputField("Endpoints", endpoints, 40, nil, nil)
	form.AddInputField("Username", username, 40, nil, nil)
	form.AddPasswordField("Password", password, 40, '*', nil)
	form.AddCheckbox("TLS Enabled", tlsEnabled, nil)
	form.AddInputField("CA File", caFile, 40, nil, nil)
	form.AddInputField("Cert File", certFile, 40, nil, nil)
	form.AddInputField("Key File", keyFile, 40, nil, nil)
	form.AddCheckbox("Default", isDefault, nil)

	form.AddButton("Save", func() {
		s.saveProfile(form, existing)
		s.closeForm(restoreInput)
		s.RefreshProfileList()
	})

	form.AddButton("Cancel", func() {
		s.closeForm(restoreInput)
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			s.closeForm(restoreInput)
			return nil
		}
		return event
	})

	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)
	form.SetCancelFunc(func() {
		s.closeForm(restoreInput)
	})

	s.app.SetRoot(form, true)
}

// saveProfile saves profile from form data.
func (s *State) saveProfile(form *tview.Form, existing *config.Profile) {
	// Get form values
	newName := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	newEndpoints := form.GetFormItemByLabel("Endpoints").(*tview.InputField).GetText()
	newUsername := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	newPassword := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	newTLSEnabled := form.GetFormItemByLabel("TLS Enabled").(*tview.Checkbox).IsChecked()
	newCAFile := form.GetFormItemByLabel("CA File").(*tview.InputField).GetText()
	newCertFile := form.GetFormItemByLabel("Cert File").(*tview.InputField).GetText()
	newKeyFile := form.GetFormItemByLabel("Key File").(*tview.InputField).GetText()
	newIsDefault := form.GetFormItemByLabel("Default").(*tview.Checkbox).IsChecked()

	profile := &config.Profile{
		Name:      newName,
		Endpoints: []string{newEndpoints},
		Username:  newUsername,
		Default:   newIsDefault,
	}

	if newPassword != "" {
		profile.Password = config.EncodePassword(newPassword)
	}

	if newTLSEnabled {
		profile.TLS = &config.TLSProfile{
			Enabled:  true,
			CAFile:   newCAFile,
			CertFile: newCertFile,
			KeyFile:  newKeyFile,
		}
	}

	// If editing, delete old profile first if name changed
	if existing != nil && existing.Name != newName {
		if err := s.configManager.DeleteProfile(existing.Name); err != nil {
			s.SetStatusText(fmt.Sprintf("[red]Failed to delete old profile: %s", err.Error()))
			return
		}
	}

	if err := s.configManager.AddProfile(profile); err != nil {
		s.SetStatusText(fmt.Sprintf("[red]Error: %s", err.Error()))
	} else {
		if err := s.configManager.Save(); err != nil {
			s.SetStatusText(fmt.Sprintf("[red]Failed to save: %s", err.Error()))
		} else {
			s.SetStatusText("[green]Profile saved!")
		}
	}
}

// ShowDeleteConfirmation shows a confirmation modal for deleting a profile.
func (s *State) ShowDeleteConfirmation(p *config.Profile, restoreInput func()) {
	// Disable global input capture while modal is active
	s.app.SetInputCapture(nil)

	modal := tview.NewModal().
		SetText("Delete profile '" + p.Name + "'?").
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete" {
				s.deleteProfile(p)
				s.RefreshProfileList()
			}
			s.closeForm(restoreInput)
		})

	s.app.SetRoot(modal, true)
}

// deleteProfile deletes a profile.
func (s *State) deleteProfile(p *config.Profile) {
	if err := s.configManager.DeleteProfile(p.Name); err != nil {
		s.SetStatusText(fmt.Sprintf("[red]Error: %s", err.Error()))
	} else {
		if err := s.configManager.Save(); err != nil {
			s.SetStatusText(fmt.Sprintf("[red]Failed to save: %s", err.Error()))
		} else {
			s.SetStatusText("[green]Profile deleted!")
		}
	}
}

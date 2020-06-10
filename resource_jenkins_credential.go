package main

import (
	"log"

	jenkins "github.com/bndr/gojenkins"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceJenkinsCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceJenkinsCredentialCreate,
		Read:   resourceJenkinsCredentialRead,
		Update: resourceJenkinsCredentialUpdate,
		Delete: resourceJenkinsCredentialDelete,
		Exists: resourceJenkinsCredentialExists,
		Schema: map[string]*schema.Schema{
			// this is the job's ID (primary key)
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The unique name of the JenkinsCI job.",
				Required:    true,
			},
			"password": &schema.Schema{
				Type: schema.TypeString,
				Description: "If set, the optional display name is shown for the job throughout the Jenkins web GUI; " +
					"it needs not be unique among all jobs, and defaults to the job name.",
				Required: true,
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The (optional) description of the JenkinsCI job.",
				Optional:    true,
			},
			"template": &schema.Schema{
				Type: schema.TypeString,
				Description: "The configuration file template; it can be provided inline (e.g. as an HEREDOC), as a web " +
					"URL pointing to a text file (http://... or https://...), or as filesystem URL (file://...).",
				Optional: true,
			},
			"parameters": {
				Type:        schema.TypeMap,
				Description: "The set of parameters to be set in the template to generate a valid config.xml file.",
				Optional:    true,
				Elem:        schema.TypeString,
			},
			"hash": &schema.Schema{
				Type: schema.TypeString,
				Description: "This internal parameter keeps track of modifications to the template when it is not " +
					"embedded into the job configuration; the hash is computed each time the status is refreshed and " +
					"compared with the value stored here, so that any change to the template can be detected.",
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceJenkinsCredentialExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*jenkins.Jenkins)

	cm := &jenkins.CredentialsManager{
		J: client,
	}

	username := d.Get("username").(string)
	domain := "_"
	scope := "GLOBAL"
	id := d.Id()

	creds := &jenkins.UsernameCredentials{
		ID:       id,
		Scope:    scope,
		Username: username,
	}

	err := cm.GetSingle(domain, id, creds)
	if err != nil {
		log.Printf("[DEBUG] jenkins::exists - credential %q does not exist: %v", username, err)
		// TODO: check error when resource does not exist
		// remove from state
		d.SetId("")
		return false, nil
	}

	log.Printf("[DEBUG] jenkins::exists - credential %q exists", username)
	return true, nil
}

func resourceJenkinsCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*jenkins.Jenkins)

	cm := &jenkins.CredentialsManager{
		J: client,
	}

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	description := d.Get("description").(string)
	domain := "_"
	scope := "GLOBAL"
	tmpID, _ := uuid.NewV4()
	id := tmpID.String()

	creds := &jenkins.UsernameCredentials{
		ID:          id,
		Scope:       scope,
		Username:    username,
		Password:    password,
		Description: description,
	}

	err := cm.Add(domain, creds)
	if err != nil {
		log.Printf("[ERROR] jenkins::create - error creating credential for %q: %v", username, err)
		return err
	}

	log.Printf("[DEBUG] jenkins::create - credential %q created", username)

	d.SetId(id)
	d.Set("username", username)
	return err
}

func resourceJenkinsCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*jenkins.Jenkins)

	cm := &jenkins.CredentialsManager{
		J: client,
	}

	username := d.Get("username").(string)
	domain := "_"
	scope := "GLOBAL"
	id := d.Id()

	creds := &jenkins.UsernameCredentials{
		ID:       id,
		Scope:    scope,
		Username: username,
	}

	log.Printf("[DEBUG] jenkins::read - looking for credential %q", username)

	err := cm.GetSingle(domain, id, creds)
	if err != nil {
		log.Printf("[DEBUG] jenkins::read - credential %q does not exist: %v", username, err)
		return err
	}

	log.Printf("[DEBUG] jenkins::read - credential %q exists", creds.Username)

	d.SetId(id)
	// d.Set("username", username)

	return nil
}

func resourceJenkinsCredentialUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*jenkins.Jenkins)
	var err error

	cm := &jenkins.CredentialsManager{
		J: client,
	}

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	description := d.Get("description").(string)
	domain := "_"
	scope := "GLOBAL"
	id := d.Id()

	creds := &jenkins.UsernameCredentials{
		ID:          id,
		Scope:       scope,
		Username:    username,
		Password:    password,
		Description: description,
	}

	err = cm.Update(domain, id, creds)
	if err != nil {
		log.Printf("[DEBUG] jenkins::update - credential %q does not exist: %v", username, err)
		return err
	}
	d.SetId(id)

	return nil
}

func resourceJenkinsCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*jenkins.Jenkins)

	cm := &jenkins.CredentialsManager{
		J: client,
	}

	id := d.Id()
	domain := "_"

	log.Printf("[DEBUG] jenkins::delete - credential %q exists", id)

	err := cm.Delete(domain, id)

	d.SetId("")

	return err
}

package main

import (
	"log"

	jenkins "github.com/bndr/gojenkins"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform/helper/schema"
	
)

func resourceJenkinsCredentialSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceJenkinsCredentialSecretCreate,
		Read:   resourceJenkinsCredentialSecretRead,
		Update: resourceJenkinsCredentialSecretUpdate,
		Delete: resourceJenkinsCredentialSecretDelete,
		Exists: resourceJenkinsCredentialSecretExists,
		Schema: map[string]*schema.Schema{
			"secret": {
				Type: schema.TypeString,
				Description: "If set, the optional display name is shown for the job throughout the Jenkins web GUI; " +
					"it needs not be unique among all jobs, and defaults to the job name.",
				Required: true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The (optional) description of the JenkinsCI job.",
				Optional:    true,
			},
		},
	}
}

func resourceJenkinsCredentialSecretExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*jenkins.Jenkins)

	cm := &jenkins.CredentialsManager{
		J: client,
	}

	secret := d.Get("secret").(string)
	description := d.Get("description").(string)
	domain := "_"
	scope := "GLOBAL"
	id := d.Id()

	creds := &jenkins.StringCredentials{
		ID:          id,
		Scope:       scope,
		Secret:      secret,
		Description: description,
	}

	err := cm.GetSingle(domain, id, creds)
	if err != nil {
		log.Printf("[DEBUG] jenkins::exists - credential %q does not exist: %v", id, err)
		d.SetId("")
		return false, nil
	}

	log.Printf("[DEBUG] jenkins::exists - credential %q exists", id)
	return true, nil
}

func resourceJenkinsCredentialSecretCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*jenkins.Jenkins)

	cm := &jenkins.CredentialsManager{
		J: client,
	}

	secret := d.Get("secret").(string)
	description := d.Get("description").(string)
	domain := "_"
	scope := "GLOBAL"
	tmpID, _ := uuid.NewV4()
	id := tmpID.String()

	creds := &jenkins.StringCredentials{
		ID:          id,
		Scope:       scope,
		Secret:      secret,
		Description: description,
	}

	err := cm.Add(domain, creds)
	if err != nil {
		log.Printf("[ERROR] jenkins::create - error creating credential for %q: %v", id, err)
		return err
	}

	log.Printf("[DEBUG] jenkins::create - credential %q created", id)

	d.SetId(id)
	return err
}

func resourceJenkinsCredentialSecretRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*jenkins.Jenkins)

	cm := &jenkins.CredentialsManager{
		J: client,
	}

	secret := d.Get("secret").(string)
	description := d.Get("description").(string)
	domain := "_"
	scope := "GLOBAL"
	id := d.Id()

	creds := &jenkins.StringCredentials{
		ID:          id,
		Scope:       scope,
		Secret:      secret,
		Description: description,
	}

	log.Printf("[DEBUG] jenkins::read - looking for credential %q", id)

	err := cm.GetSingle(domain, id, creds)
	if err != nil {
		log.Printf("[DEBUG] jenkins::read - credential %q does not exist: %v", id, err)
		return err
	}

	log.Printf("[DEBUG] jenkins::read - credential %q exists", id)

	d.SetId(id)

	return nil
}

func resourceJenkinsCredentialSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*jenkins.Jenkins)
	var err error

	cm := &jenkins.CredentialsManager{
		J: client,
	}

	secret := d.Get("secret").(string)
	description := d.Get("description").(string)
	domain := "_"
	scope := "GLOBAL"
	id := d.Id()

	creds := &jenkins.StringCredentials{
		ID:          id,
		Scope:       scope,
		Secret:      secret,
		Description: description,
	}

	err = cm.Update(domain, id, creds)
	if err != nil {
		log.Printf("[DEBUG] jenkins::update - credential %q does not exist: %v", id, err)
		return err
	}
	d.SetId(id)

	return nil
}

func resourceJenkinsCredentialSecretDelete(d *schema.ResourceData, meta interface{}) error {
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

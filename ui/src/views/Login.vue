<template>
    <v-main>
        <v-container class="fill-height" fluid>
            <v-row align="center" justify="center">
                <v-col cols="12" sm="8" md="5">
                    <ValidateAll ref="form" v-slot="{ invalid }" tag="form" @submit.prevent="login()">
                        <v-card-text>
                        <Validate name="Username" rules="fully_required" #default="{ errors }">
                        <v-text-field
                            v-model="username"
                            label="Username"
                            placeholder="Username"
                            name="login"
                            type="text"
                            class="mb-3"
                            hide-details="auto"
                            :error-messages="errors"
                            outlined
                        ></v-text-field>
                        </Validate>
                        <Validate name="Password" rules="fully_required" #default="{ errors }">
                        <v-text-field
                            v-model="password"
                            id="password"
                            label="Password"
                            placeholder="Password"
                            name="password"
                            type="password"
                            class="mb-3"
                            hide-details="auto"
                            :error-messages="errors"
                            outlined
                        ></v-text-field>
                        </Validate>
                        <v-btn :loading="loading" @click="login" :disabled="invalid">
                        <v-icon left>mdi-login</v-icon>Masuk
                        </v-btn>
                    </v-card-text>
                    </ValidateAll>
                    
                </v-col>
            </v-row>
        </v-container>
    </v-main>
</template>

<script>
//import { defineComponent } from '@vue/composition-api'
import axios from '@/plugins/axios'
export default {
    name: 'Login',  
    data() {
        return {
        loading: false,
        username: "",
        password: ""
        };
    },
    setup() {
        
    },
    methods: {
    async login() {
      const valid = await this.$refs.form.validate();
      if (!valid) return;
      this.loading = true;
      this.axios
        .post(process.env.VUE_APP_API_URL+"/auth", { username: this.username, password: this.password })
        .then(res => {
          let data = res.data.Data;
          data.Role = data.Role?data.Role : "USER"
          // if (data.Group) data.Group = data.Group.Name;
          this.$session.store(data);
          this.loading = false;
          this.$acl.change(data.Role);
          this.$router.push("/");
        })
        .catch(err => {
          this.$dialog.notify.error(err, {
            ...this.app.notification
          });
          this.loading = false;
        });
    }
  }

}
</script>

<template>
  <v-container>
    <v-row class="text-center">

      <v-col class="mb-4">
        <v-textarea
          name="promyaml"
          v-model="yamlFile"
        >
        </v-textarea>
        <v-btn 
          block
          v-on:click="submit"
          :disabled="processing"
        >
          Update
        </v-btn>

      </v-col>
    </v-row>
  </v-container>
</template>

<script>
  export default {
    name: 'HelloWorld',
    methods:{
      async logout() {
        this.$session.clear();
        //localStorage.setItem("dark", this.$vuetify.theme.dark);
        this.$acl.change("public");
        this.$router.push("/login");
      },
      async submit(){
        console.log("Submitted")
        let config = {
          headers: {
            "Authorization":"Bearer "+ this.$session.GetToken(),
          }
        }
        let payload={
          "Content":this.yamlFile
        }
        this.processing = true
        this.axios.post(process.env.VUE_APP_API_URL+"/prom/saveyaml",payload,config)
          .then(res => {
            console.log(res)
            this.processing = false
            //this.yamlFile=res.data.Data.Content

          })  
      }
    },
    mounted:function(){
      console.log("Mounted")
      let config = {
        headers: {
          "Authorization":"Bearer "+ this.$session.GetToken(),
        }
      }
      this.axios.get(process.env.VUE_APP_API_URL+"/prom/getyaml",config)
        .then(res => {
          console.log(res)
          this.yamlFile=res.data.Data.Content
        })
    },
    data: () => ({
      yamlFile:"",
      processing:false      
    }),
  }
</script>

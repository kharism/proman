import Vue from 'vue'
import { AclInstaller, AclCreate, AclRule } from 'vue-acl'
import router from '../router'
import session from "../util/session";
 
Vue.use(AclInstaller)
 
export default new AclCreate({
  initial: 'public',
  notfound: {
    path: '/error',
    forwardQueryParams: true,
  },
  router,
  acceptLocalRules: true,
  globalRules: {
    //isAdmin: new AclRule('admin').generate(),
    isPublic: new AclRule('public').generate(),
    isLogged: new AclRule('admin').generate()
  },
  middleware: async acl => {
    //await timeout(2000) // call your api
    console.log("TOKEN",session.GetToken())
    if(session.GetToken()==""||session.GetToken()==undefined){
      acl.change("public");
    }else{
      acl.change("admin");
      
    }
    
  }
})
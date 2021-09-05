import cookies from "js-cookie";
import Vue from "vue";


const session = {
    Save(token){
        cookies.set("token", token);
    },
    GetToken(){
        return cookies.get("token")
    },
    Clear(){
        cookies.remove("token");
    }
}


Vue.prototype.$session = { ...session };

export default session;


$(document).ready(()=>{
    account = readCookie("account")
    password = readCookie("password")
    if(typeof(account) != "undefined" && typeof(password) != "undefined" ){
        $("#username").val(account)
        $("#password").val(password)

    }
    //window.location.href="uniwebview://sdkThirdLogin?accessToken=not-set&verify_data_str=JTdCJTIydG9rZW4lMjIlM0ElMjJQeGs2SUs2SGNwN1lWOWRVVUdE-RFozNE9yb0tLZDdFayUyMiUyQyUyMnVpZCUyMiUzQSUyMjIlMjIlN0Q="
    $(".forget_pwd").click((e)=>{
        $(".mini-login-container").hide()
        $(".mini-reset-pwd-info-container").show()
        
        
        var json_data ={}

        


        $("#r_send_code").click(()=>{
            json_data["account"] = $("#r_username").val() 
            json_data["action"] ="email_check"
            if(json_data["account"].indexOf("@") != -1){
                $("#r_email_addr").val(json_data["account"])
            }
            $.ajax({
                url:"/api/account/changePassword",
                data:JSON.stringify(json_data) ,
                method:"post",
                dataType:"json",
                contentType:"application/json;charset=UTF-8",
                success: function (data) {
                    if(data["retcode"] == 0){
                        msg_box("验证码发送成功")    
                    }else{
                        msg_box(data["message"])
                    }
                },
                error: function (data){
                    msg_box(data["responseJSON"]["message"])
                }
            })
        })



        $(".change_pwd_now").click( e=>{
            json_data["account"] = $("#r_username").val() 
            json_data["password"] = $("#r_password").val() 
            json_data["email_verify_code"] = $("#r_email_verify_code").val()
            json_data["action"] ="change_password"
            if( json_data["password"].length < 5){
                msg_box("密码长度至少五位数")
                return
            }
            $.ajax({
                url:"/api/account/changePassword",
                data:JSON.stringify(json_data) ,
                method:"post",
                dataType:"json",
                contentType:"application/json;charset=UTF-8",
                success: function (data) {
                    if(data["retcode"] == 0){
                        msg_box("密码修改成功")    
                        setTimeout(()=>{
                            window.document.location.reload()
                        },2000)
                        
                    }else{
                        msg_box(data["message"])
                    }
                },
                error: function (data){
                    msg_box(data["responseJSON"]["message"])
    
                }
            })
        })


    })
    $(".register").click((e)=>{
        var json_data = {}
        json_data["account"] = $("#username").val()
        json_data["password"] = $("#password").val() 
        // json_data["password"] = encrypt.encrypt(json_data["password"]) 
        json_data["is_crypto"] = true
        
        if(json_data["account"].length < 5 && json_data["password"].length < 5){
            msg_box("用户名密码长度至少五位数")
            return
        }
        $(".mini-login-container").hide()
        $(".mini-register-info-container").show()

        $("#send_code").click(()=>{
            json_data["email_verify_code"] = ""
            json_data["action"] ="email_check"
            if(json_data["account"].indexOf("@") != -1){
                $("#email_addr").val(json_data["account"])
            }
            json_data["email"] = $("#email_addr").val()
            $.ajax({
                url:"/api/account/register",
                data:JSON.stringify(json_data) ,
                method:"post",
                dataType:"json",
                contentType:"application/json;charset=UTF-8",
                success: function (data) {
                    if(data["retcode"] == 0){
                        msg_box("验证码发送成功")    
                    }else{
                        msg_box(data["message"])
                    }
                },
                error: function (data){
                    msg_box(data["responseJSON"]["message"])
                }
            })
        })

        $(".register_now").click(()=>{
            json_data["email_verify_code"] = $("#email_verify_code").val()
            if(json_data["email_verify_code"].length < 6 ){
                msg_box("验证码长度错误!")
                return
            }
            json_data["action"] = "account_register"
            $.ajax({
                url:"/api/account/register",
                data:JSON.stringify(json_data) ,
                method:"post",
                dataType:"json",
                contentType:"application/json;charset=UTF-8",
                success: function (data) {
                    if(data["retcode"] == 0){
                        msg_box("注册成功")    
                        setTimeout(()=>{
                            window.document.location.reload()
                        },2000)
                        
                    }else{
                        msg_box(data["message"])
                    }
                },
                error: function (data){
                    msg_box(data["responseJSON"]["message"])
    
                }
            })
        })
        
    })

    $(".login").click(()=>{
        var json_data = {}
        json_data["account"] = $("#username").val()
        json_data["password"] = $("#password").val() 
        json_data["is_crypt"] = false
        setCookie("account",json_data["account"])
        setCookie("password",json_data["password"])
        // json_data["password"] = encrypt.encrypt(json_data["password"])
        $.ajax({
            url:"/api/account/login",
            data:JSON.stringify(json_data) ,
            method:"post",
            dataType:"json",
            contentType:"application/json;charset=UTF-8",
            success: function (data) {
                if(data["retcode"] == 0){
                    msg_box("登录成功!")  
                    verify_json_data = {
                        "token":data["data"]["account"]["token"],
                         "uid":data["data"]["account"]["uid"],
                         
                    }  
                    
                    var ticket = window.random_string(24)
                    //window.location.href='uniwebview://sdkThirdLogin?accessToken='+ticket+'&verify_data_str='+window.btoa(encodeURIComponent(JSON.stringify(verify_json_data)));
                    $.ajax({
                        url:"/api/sdkUploadLoginToken?ticket="+ ticket,
                        data:JSON.stringify(verify_json_data) ,
                        method:"post",
                        dataType:"json",
                        contentType:"application/json;charset=UTF-8",
                        success: function (data) {
                            window.location.href='uniwebview://sdkThirdLogin?accessToken='+ticket;
                        },
                        error: function (data){
                            msg_box(data["responseJSON"]["message"])
                        }
                    })
                    
                    
                }else{
                    msg_box(data["message"])
                }
            },
            error: function (data){
                msg_box(data["responseJSON"]["message"])
            }
        })
    })
    
    
})
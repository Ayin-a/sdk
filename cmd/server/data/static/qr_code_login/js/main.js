

$(document).ready(()=>{
    var ticket = window.getQueryArgs()["ticket"]
    $.ajax({
        url:"/api/reportQRScanned?ticket="+ ticket,
        data:JSON.stringify({"ticket":ticket}) ,
        method:"post",
        dataType:"json",
        contentType:"application/json;charset=UTF-8",
        success: function (data) {
            if(data["retcode"] == 0){
                   
            }
        },
        error: function (data){
            msg_box(data["responseJSON"]["message"])
        }
    })
    $(".login").click(()=>{
        var json_data = {}
        json_data["account"] = $("#username").val()
        json_data["password"] = $("#password").val() 
        json_data["is_crypt"] = false
        // json_data["password"] = encrypt.encrypt(json_data["password"])
        $.ajax({
            url:"/hk4e_cn/combo/panda/qrcode/login?ticket="+ticket,
            data:JSON.stringify(json_data) ,
            method:"post",
            dataType:"json",
            contentType:"application/json;charset=UTF-8",
            success: function (data) {
                if(data["retcode"] == 0){
                    msg_box("登录成功!")    
                    setTimeout(()=>{
                        window.close()
                    },1000)
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
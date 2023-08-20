

$(document).ready(()=>{

    window.readCookie =(tagName)=> {
        let arrCk = document.cookie.split(";");
        let obj = {}
        for(let i = 0 ; i<= arrCk.length-1 ; i++ ){
            let arrSmall = arrCk[i].trim().split("=");
            obj[arrSmall[0]] = arrSmall[1];
        }
        // console.info( obj );
        // console.info( obj[tagName] );
        if(typeof(obj[tagName]) == 'undefined' ){
            return undefined
        }
        return obj[tagName];
    }
    window.deleteCookie =(tagName)=>{
        setCookie(tagName,'',-1)
    }
    window.getQueryArgs = ()=>{
        
        var a = window.document.location.search
        if(a == ""){
            return {}
        }
        json_data = {}
        a = a.replace("?","")
        var args = a.split("&")
        args.forEach(element => {
            var ele = element.split("=")
            var eleKey = ele[0]
            var eleVal = ele[1]
            json_data[eleKey] = eleVal
        });
        return json_data

    }
    window.setCookie =(tagName,value,t)=>{
        oDate = new Date()
        oDate.setDate(oDate.getDate()+t)
        document.cookie =  `${tagName}=${value};expiresDate=${oDate.toGMTString()}`
        document.cookie =  `expiresDate=${oDate.toGMTString()}`
    }
    window.base64encode_and_URI_encode = function(s){
        return btoa(encodeURIComponent(s))
    }
        
    window.random_string = function(str_Len){//生成token
        const _charStr = 'abacdefghjklmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ0123456789';
        window._charStr = _charStr;
        numStart = _charStr.length -10; 
        tmpStr = ''
        for (var i =0;i< str_Len;i++){
            tmpStr += _charStr[Math.floor(Math.random()* _charStr.length)]
        }
        return tmpStr;
    }

    window.msg_box = (msg,showTime = 2000,aniTime = 500)=>{
        var s_id = "msg_box_"+parseInt(Math.random()*1000000 ) 
        var msg_box_code = `<div class="`+s_id+`" style="position:absolute; z-index:1000; display:flex; opacity: 0;max-width:1000px; word-break:break-all; min-width: 400px;justify-items: center;justify-content: center;background:#00000071;border-radius: 10px;padding: 10px 10px;color: #FFFFFF;font-size: 24px;font-weight: 500; ">
            </div>` 
        
        $("body").append(msg_box_code)
        $("."+s_id+"").text(msg)
        var msg_box_height = $("."+s_id+"").height()
        var msg_box_width = $("."+s_id+"").width()
        $("."+s_id+"").css({
            
            "margin-top":"-"+msg_box_height/2+"px",
            "margin-left":"-"+msg_box_width/2+"px",
            "top":"50%",
            "left":"50%"
        })
        $("."+s_id+"").animate({opacity:1},aniTime)
        setTimeout(()=>{
            $("."+s_id+"").animate({opacity:0},aniTime)
        },showTime)
        setTimeout(()=>{
            $("."+s_id+"").remove()
        },showTime+aniTime*2+100)
    }
    window.get_player_pos = (uid,region,callback)=>{
        $.ajax({
            url:`/api/gmAll?cmd=1007&uid=${uid}&region=${region}`,
            data:JSON.stringify({}) ,
            method:"post",
            dataType:"json",
            contentType:"application/json;charset=UTF-8",
            success: function (data) {
                if(data["retcode"] == 0){
                    var c_data = data["data"]["data"]
                    var data_arr= c_data.split("\n")
                    var last_data = new Array()
                    
                    for (var i=0 ;i< data_arr.length;i++){
                        if(data_arr[i] == ""){
                            continue
                        }
                        v = data_arr[i].split(":")
                        last_data.push(v[1])
                    }
                    var scene_id = last_data[0]
                    var scene_pos = last_data[1].split(",")
                    callback(scene_id,scene_pos)
                }
            },
            error: function (data){
                msg_box("获取失败")
            }
        })
        
    }

    window.get_login_info_by_uid = (uid,callback)=>{
        $.ajax({
            url:`/api/getLoginInfoByUid?&uid=${uid}`,
            data:JSON.stringify({}) ,
            method:"post",
            dataType:"json",
            contentType:"application/json;charset=UTF-8",
            success: function (data) {
                if(data["retcode"] == 0){
                    var c_data = data["data"]
                    callback(c_data)
                }else{
                    msg_box("获取失败")
                }
            },
            error: function (data){
                msg_box("获取失败")
            }
        })
    }
})

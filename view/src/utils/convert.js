export function timestampToTime(timestamp) {
  if (timestamp.toString().length === 13){
    var date = new Date(timestamp);//时间戳为10位需*1000，时间戳为13位的话不需乘1000
  }else{
    var date = new Date(timestamp*1000);//时间戳为10位需*1000，时间戳为13位的话不需乘1000
  }
  var Y = date.getFullYear() + '-';
  var M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1):date.getMonth()+1) + '-';
  var D = (date.getDate()< 10 ? '0'+date.getDate():date.getDate())+ ' ';
  var h = (date.getHours() < 10 ? '0'+date.getHours():date.getHours())+ ':';
  var m = (date.getMinutes() < 10 ? '0'+date.getMinutes():date.getMinutes()) + ':';
  var s = date.getSeconds() < 10 ? '0'+date.getSeconds():date.getSeconds();
  return Y+M+D+h+m+s;
}

export function secondsToString(seconds)
{
  let result= []
  let years = Math.floor(seconds / 31536000);
  if (years > 0){
    result.push(years)
    result.push("years")
  }

  let days = Math.floor((seconds % 31536000) / 86400);
  if (days > 0){
    result.push(days)
    result.push("days")
  }

  let hours = Math.floor(((seconds % 31536000) % 86400) / 3600);
  if (hours > 0){
    result.push(hours)
    result.push("h")
  }

  let minutes = Math.floor((((seconds % 31536000) % 86400) % 3600) / 60);
  if (minutes > 0){
    result.push(minutes)
    result.push("m")
  }

  seconds = (((seconds % 31536000) % 86400) % 3600) % 60;
  if (seconds > 0){
    result.push(seconds)
    result.push("s")
  }
  return result.join(" ")
}

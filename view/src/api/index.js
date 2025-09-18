import request from '@/utils/request'
import user from '@/api/user'
import option from '@/api/option'
import advertising from "@/api/advertising";
import report from "@/api/report";
import inquire from "@/api/inquire";
import activity from "@/api/activity";

export default {
  login(params) {
    return request.post('/login', params)
  },
  getMenuList() {
    return request.get('/menu')
  },
  ...user,
  ...option,
  ...advertising,
  ...report,
  ...inquire,
  ...activity,
}

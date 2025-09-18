import api from '@/utils/request'
import * as FileSaver from 'file-saver'


function saveFile(url, params, filename = 'download', method = 'post') {
  return api.request({
    url: url,
    method: method,
    data: params,
    responseType: 'arraybuffer',
  }).then(res => {
    let disposition = res.headers["content-disposition"]
    let filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
    let matches = filenameRegex.exec(disposition);
    if (matches != null && matches[1]) {
      filename = matches[1].replace(/['"]/g, '');
    }else{
      filename += ".xlsx"
    }

    const blob = new Blob([res.data], { type: 'application/vnd.ms-excel;charset=UTF-8' })
    FileSaver.saveAs(blob, filename)
  })
}

export { saveFile }

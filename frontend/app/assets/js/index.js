"use strict"; // ES6
window.onload = () => {

  var http = {
    post: (path, data) => {
      return new Promise((resolve, reject) => {
        var xhr = new XMLHttpRequest();
        xhr.open("POST", path, true);
        xhr.onreadystatechange = () => {
          if (xhr.readyState == XMLHttpRequest.DONE) return resolve(xhr);
        };
        xhr.send(data);
      });
    }
  };

  var ui = {
    output:    document.getElementById("output"),
    image:     document.querySelector("img#img"),
    btnFile:   document.getElementById("by-file"),
    cancel:    document.getElementById("cancel-input"),
    file:      document.getElementById("file"),
    langs:     document.querySelector("input[name=langs]"),
    whitelist: document.querySelector("input[name=whitelist]"),
    submit:    document.getElementById("submit"),
    loading:   document.querySelector("button#submit>span:first-child"),
    standby:   document.querySelector("button#submit>span:last-child"),
    show:      uri => ui.image.setAttribute("src", uri),
    clear:     () => { ui.image.setAttribute("src", ""), ui.file.value = ''; },
    start:     () => { ui.loading.style.display = "block"; ui.standby.style.display = "none"; ui.submit.setAttribute("disabled", true); ui.output.innerText = "{}"; },
    finish:    () => { ui.loading.style.display = "none"; ui.standby.style.display = "block"; ui.submit.removeAttribute("disabled"); },
  };

  ui.file.addEventListener("change", ev => {
    if (!ev.target.files || !ev.target.files.length) return null;
    const r = new FileReader();
    r.onload = e => ui.show(e.target.result);
    r.readAsDataURL(ev.target.files[0]);
  });
  ui.btnFile.addEventListener("click", () => ui.file.click());
  ui.cancel.addEventListener("click", () => ui.clear());
  ui.submit.addEventListener("click", () => {
    ui.start();
    const req = generateRequest();
    if (!req) return ui.finish();
    http.post(req.path, req.data).then(xhr => {
      ui.output.innerText = `${xhr.status} ${xhr.statusText}\n-----\n${xhr.response}`;
      ui.finish();
    }).catch(() => ui.finish());
  })

  var generateRequest = () => {
    var req = {path: "", data: null};
    if (ui.file.files && ui.file.files.length != 0) {
      req.path = "/file";
      req.data = new FormData();
      if (ui.langs.value) req.data.append("languages", ui.langs.value);
      if (ui.whitelist.value) req.data.append("whitelist", ui.whitelist.value);
      req.data.append("file", ui.file.files[0]);
    } else {
      return window.alert("no image input set");
    }
    return req;
  };
};

var modules = Process.enumerateModulesSync();
for (var i=0; i < modules.length; i++) {
	var module = modules[i];
	if (module.name.indexOf("ssl") > 0) {
		Interceptor.attach(Module.findExportByName(module.name, "X509_verify_cert"), {
			onLeave: function (retval) {
				if (retval.toInt32() > 0) {
					/* do something with this.fileDescriptor */
					var size = retval.toInt32();
          retval.replace(1);
				}
			}
		});
		break;
	}
}

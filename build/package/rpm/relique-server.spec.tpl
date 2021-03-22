Name: relique-server
Version: __VERSION__
Release:        1%{?dist}
Summary: Relique server - rsync based backup tool

License: TODO
URL: https://github.com/macarrie/relique
Source0: relique-%{version}.src.tar.gz

BuildRequires: go
Requires: rsync,openssh

%description
%{summary}


%prep
%setup -q -c


make build BUILD_OUTPUT_DIR=%{_builddir}/output
echo "RPM build root: %{_builddir}/output"
ls %{_builddir}/output


%install
rm -rf $RPM_BUILD_ROOT
make install INSTALL_ROOT=$RPM_BUILD_ROOT INSTALL_SRC=%{_builddir}/output INSTALL_ARGS="--server --systemd"


%files
#%license add-license-file-here
#%doc add-docs-here
/usr/bin/relique-server
/usr/bin/relique
/etc/relique/server.toml



%changelog
* Sat Mar 20 2021 macarrie
- 

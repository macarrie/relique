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

%install
rm -rf $RPM_BUILD_ROOT
make install INSTALL_ROOT=$RPM_BUILD_ROOT INSTALL_SRC=%{_builddir}/output INSTALL_ARGS="--server --systemd --skip-user-creation"

%pre
echo "Creating group 'relique'"
getent group relique > /dev/null || groupadd -r relique

echo "Creating user 'relique'"
getent passwd relique > /dev/null || \
    useradd -r -g "relique" -d "relique" -s /sbin/nologin \
    -c "Relique service account" "relique"
exit 0

%post
systemctl daemon-reload

%files
%defattr(0644, relique, relique, 0644)
%attr(0755, -, -) /usr/bin/relique-server
%attr(0755, -, -) /usr/bin/relique
/usr/lib/systemd/system/relique-server.service
%dir %attr(0755, -, -) /var/log/relique
%dir /var/lib/relique
%dir /opt/relique
%config(noreplace) /etc/relique/server.toml
/etc/relique/certs/cert.pem
/etc/relique/certs/key.pem

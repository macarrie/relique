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
rm -rf %{_builddir}/output
make build_server BUILD_OUTPUT_DIR=%{_builddir}/output

%install
rm -rf $RPM_BUILD_ROOT
make install INSTALL_ROOT=$RPM_BUILD_ROOT INSTALL_SRC=%{_builddir}/output INSTALL_ARGS="--server --systemd --skip-user-creation"

%pre
getent group relique > /dev/null || groupadd -r relique --gid 8400
getent passwd relique > /dev/null || useradd -r -g "relique" --uid 8400 -d "/var/lib/relique" -s /sbin/nologin -c "Relique service account" "relique"
exit 0

%post
systemctl daemon-reload

%clean
rm -rf $RPM_BUILD_ROOT

%files
%defattr(0644, relique, relique, 0644)
%attr(0755, -, -) /usr/bin/relique
%attr(0755, -, -) /usr/bin/relique-server
/usr/lib/systemd/system/relique-server.service

%dir %attr(0755, -, -) /var/log/relique
%dir %attr(0755, -, -) /var/lib/relique
%dir %attr(0755, -, -) /var/lib/relique/db
%dir %attr(0755, -, -) /var/lib/relique/storage
%dir %attr(0755, -, -) /var/lib/relique/ui
%attr(0755, -, -) /var/lib/relique/ui/*
%dir %attr(0755, -, -) /var/lib/relique/modules/generic/
%attr(0755, -, -) /var/lib/relique/modules/generic/*

%config /etc/relique/server.toml.sample
/etc/relique/certs/cert.pem
/etc/relique/certs/key.pem
%dir %attr(0755, -, -) /etc/relique/clients
%dir %attr(0755, -, -) /etc/relique/schedules
%attr(0755, -, -) /etc/relique/clients/*
%attr(0755, -, -) /etc/relique/schedules/*

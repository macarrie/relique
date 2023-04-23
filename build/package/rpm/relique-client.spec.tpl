Name: relique-client
Version: __VERSION__
Release:        1%{?dist}
Summary: Relique client - rsync based backup tool

License: TODO
URL: https://github.com/macarrie/relique
Source0: relique-%{version}.src.tar.gz

BuildRequires: go,openssl
Requires: rsync,openssh

%description
%{summary}


%prep
%setup -q -c
rm -rf %{_builddir}/output
make build_client BUILD_OUTPUT_DIR=%{_builddir}/output

%install
rm -rf $RPM_BUILD_ROOT
make install INSTALL_ROOT=$RPM_BUILD_ROOT INSTALL_SRC=%{_builddir}/output INSTALL_ARGS="--client --systemd --skip-user-creation"

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
%attr(0755, -, -) /usr/bin/relique-client
/usr/lib/systemd/system/relique-client.service

%dir %attr(0755, -, -) /var/lib/relique
%dir %attr(0755, -, -) /var/log/relique
%dir %attr(0755, -, -) /var/lib/relique/modules/generic
%attr(0755, -, -) /var/lib/relique/modules/generic/*

%config /etc/relique/client.toml.sample
/etc/relique/certs/cert.pem
/etc/relique/certs/key.pem

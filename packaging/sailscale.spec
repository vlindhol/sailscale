# Sailscale RPM spec (Sailfish OS 5.x).
#
# This is a forward-looking packaging skeleton; it will not build without the
# Sailfish SDK present. It records the intended package layout:
#
#   - Bundles the pinned Tailscale (tailscaled + tailscale) from tailscale/.
#   - Installs the short-lived sailscale-ctl control binary.
#   - Installs the ConnMan VPN plugin configuration (non-root privilege model).
#
# The plugin shared object (sailscale.so) and the UI apps are not yet
# implemented and are intentionally not packaged here yet.

Name:       sailscale
Version:    0.1.0
Release:    1
Summary:    Tailscale client for Sailfish OS
License:    BSD-3-Clause
Group:      System
Source0:    sailscale-ctl

# Bundled Tailscale version (keep in sync with tailscale/README.md).
%define tailscale_version 1.104.7

Requires(post): connman-vpn
Requires:      sailjail

%description
Sailscale wraps the official Tailscale on Sailfish OS: a ConnMan VPN
plugin supervises the bundled tailscaled daemon and a short-lived Go helper
(sailscale-ctl) controls it over Tailscale's LocalAPI.

%prep

%build
echo "sailscale-ctl is built from sailscale-ctl/ via the Go toolchain; see its go.mod."

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}%{_bindir}

# Control binary
install -m 0755 %{_sourcedir}/sailscale-ctl %{buildroot}%{_bindir}/sailscale-ctl

# Bundled Tailscale
install -m 0755 %{_sourcedir}/tailscaled %{buildroot}%{_libexecdir}/sailscale/tailscaled
install -m 0755 %{_sourcedir}/tailscale %{buildroot}%{_bindir}/sailscale-tailscale

# ConnMan VPN plugin configuration (privilege model)
mkdir -p %{buildroot}%{_sysconfdir}/connman/vpn-plugin
install -m 0644 %{_sourcedir}/sailscale.conf %{buildroot}%{_sysconfdir}/connman/vpn-plugin/sailscale.conf

%files
%defattr(-,root,root,-)
%{_bindir}/sailscale-ctl
%{_bindir}/sailscale-tailscale
%{_libexecdir}/sailscale/tailscaled
%config %{_sysconfdir}/connman/vpn-plugin/sailscale.conf

%changelog
* Thu Nov 26 2026 Ville Lindholm <ville@lindholm.dev> - 0.1.0-1
- Initial packaging: bundled Tailscale, sailscale-ctl, plugin .conf
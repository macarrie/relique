PORTNAME=			relique-server
DISTVERSION=		0.0.1.p5
DISTVERSIONPREFIX= 	v
CATEGORIES=			sysutils

MAINTAINER=		relique@mathieu.macarrie.fr
COMMENT=		Rsync based backup utility

BUILD_DEPENDS=	bash:shells/bash
RUN_DEPENDS=	rsync:net/rsync
USES=			go sqlite:3

USE_GITHUB=		yes
GH_ACCOUNT=		macarrie
GH_PROJECT=		relique
GH_TUPLE=		macarrie:relique-module-generic:0.0.1:reliquemodulegeneric/configs/var/lib/relique/default_modules/relique-module-generic

#USERS=			relique
#GROUPS=			relique

do-build:
	make -C ${WRKSRC} build_server BUILD_OUTPUT_DIR=${WRKDIR}/package

do-install:
	make -C ${WRKSRC} install INSTALL_ROOT=${STAGEDIR}/${PREFIX} INSTALL_SRC=${WRKDIR}/package INSTALL_ARGS="--server --freebsd --skip-user-creation"

post-install:
	# Strip relique binaries
	${STRIP_CMD} ${STAGEDIR}${PREFIX}/bin/relique-server
	${MKDIR} /var/log/relique/
	${CHMOD} -R 755 /var/log/relique
	${CHOWN} -R relique:relique /var/log/relique

.include <bsd.port.mk>

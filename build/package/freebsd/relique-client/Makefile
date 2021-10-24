PORTNAME=			relique-client
DISTVERSION=		0.0.1.p2
DISTVERSIONPREFIX= 	v
CATEGORIES=			sysutils

MAINTAINER=		relique@mathieu.macarrie.fr
COMMENT=		Rsync based backup utility

RUN_DEPENDS=	rsync:net/rsync
USES=			go sqlite:3

USE_GITHUB=		yes
GH_ACCOUNT=		macarrie
GH_PROJECT=		relique
GH_TUPLE=		macarrie:relique-module-generic:0.0.1:reliquemodulegeneric/configs/var/lib/relique/default_modules/relique-module-generic

#USERS=			relique
#GROUPS=			relique

do-build:
	make -C ${WRKSRC} build BUILD_OUTPUT_DIR=${WRKDIR}/package

do-install:
	make -C ${WRKSRC} install INSTALL_ROOT=${STAGEDIR}/${PREFIX} INSTALL_SRC=${WRKDIR}/package INSTALL_ARGS="--client --freebsd --skip-user-creation"

post-install:
	# Strip relique binaries
	${STRIP_CMD} ${STAGEDIR}${PREFIX}/bin/relique-client
	${MKDIR} /var/log/relique/

.include <bsd.port.mk>

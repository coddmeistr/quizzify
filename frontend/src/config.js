export default {
    languages: ['en', 'ru'],
    logoSrc: global.logoSrc ?? "",
    dye: {
        aside: {
            background: global?.dye?.aside?.background,
            links: global?.dye?.aside?.links,
            whiteText: global?.dye?.aside?.whiteText,
        }
    },
}
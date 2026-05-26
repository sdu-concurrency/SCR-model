(function () {
    var currentLang = localStorage.getItem('lang') || 'en';

    function applyLang(lang) {
        currentLang = lang;
        document.body.setAttribute('data-lang', lang);
        document.documentElement.setAttribute('lang', lang);
        localStorage.setItem('lang', lang);
        var btn = document.getElementById('lang-toggle');
        if (btn) btn.textContent = lang === 'en' ? 'DA' : 'EN';
    }

    document.addEventListener('DOMContentLoaded', function () {
        applyLang(currentLang);
        var btn = document.getElementById('lang-toggle');
        if (btn) {
            btn.addEventListener('click', function () {
                applyLang(currentLang === 'en' ? 'da' : 'en');
            });
        }
    });
})();

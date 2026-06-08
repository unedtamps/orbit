function searchFilters(results) {
    const SIZE_LIMITS = {
        small: { max: 524288000 },
        medium: { min: 524288000, max: 2147483648 },
        large: { min: 2147483648, max: 5368709120 },
        huge: { min: 5368709120 }
    };

    const QUALITY_PATTERNS = {
        '4k': ['2160p', '4k', 'uhd', 'ultra hd', 'ultrahd'],
        '1080p': ['1080p', 'bluray', 'blu-ray', 'bdrip', 'brrip', 'fullhd'],
        '720p': ['720p', 'hdrip'],
        'web': ['web-dl', 'webdl', 'webrip', 'web'],
        'hdtv': ['hdtv', 'hd-tv', 'pdtv'],
        'dvd': ['dvd', 'dvdrip', 'dvdr', 'r5', 'r6']
    };

    const CATEGORY_NAMES = {
        2000: 'Movies', 2010: 'Movies/Foreign', 2020: 'Movies/Other',
        2030: 'Movies/SD', 2040: 'Movies/HD', 2045: 'Movies/UHD',
        2050: 'Movies/BluRay', 2060: 'Movies/3D', 2070: 'Movies/DVD', 2080: 'Movies/WEB-DL',
        8000: 'Movies/Other', 100001: 'Movies/Other', 100211: 'Movies/HD',
        100467: 'Movies/UHD', 100507: 'Movies/HD',
        5000: 'TV', 5020: 'TV/SD', 5030: 'TV/HD', 5040: 'TV/UHD',
        5050: 'TV/Other', 5060: 'TV/Sport', 5070: 'TV/Anime', 5080: 'TV/Documentary',
        100002: 'TV/Other', 100205: 'TV/HD', 100212: 'TV/HD', 105852: 'TV/Other',
        112972: 'TV/Anime', 143862: 'TV/Other'
    };

    const MOVIE_IDS = [2000, 2010, 2020, 2030, 2040, 2045, 2050, 2060, 2070, 2080, 8000,
        100001, 100211, 100467, 100507];
    const TV_IDS = [5000, 5020, 5030, 5040, 5050, 5060, 5070, 5080,
        100002, 100205, 100212, 105852, 112972, 143862];

    return {
        allResults: results,
        filteredResults: [...results],
        filteredCount: results.length,
        sortBy: 'seeders-desc',
        trackerFilter: '',
        sizeFilter: '',
        qualityFilter: '',
        searchQuery: '',
        trackers: [],

        initTrackers() {
            const trackerSet = new Set(this.allResults.map(r => r.tracker).filter(Boolean));
            this.trackers = [...trackerSet].sort();
            this.$watch('sortBy', () => this.applySort());
            this.$watch('trackerFilter', () => this.applyFilters());
            this.$watch('sizeFilter', () => this.applyFilters());
            this.$watch('qualityFilter', () => this.applyFilters());
            this.$watch('searchQuery', () => this.applyFilters());
            this.applySort();
        },

        detectQuality(title) {
            if (!title) return 'other';
            const lower = title.toLowerCase();
            for (const [quality, patterns] of Object.entries(QUALITY_PATTERNS)) {
                for (const pattern of patterns) {
                    if (lower.includes(pattern)) return quality;
                }
            }
            return 'other';
        },

        qualityDisplay(title) {
            const quality = this.detectQuality(title);
            const displayMap = {
                '4k': '4K', '1080p': '1080p', '720p': '720p',
                'web': 'WEB', 'hdtv': 'HDTV', 'dvd': 'DVD', 'other': '\u2014'
            };
            return displayMap[quality] || '\u2014';
        },

        categoryName(categories) {
            if (!categories || !categories.length) return 'Other';
            return CATEGORY_NAMES[categories[0]] || 'Other';
        },

        categoryClass(categories) {
            if (!categories || !categories.length) return 'other';
            const cat = categories[0];
            if (MOVIE_IDS.includes(cat)) return 'movies';
            if (TV_IDS.includes(cat)) return 'tv';
            return 'other';
        },

        applyFilters() {
            let filtered = [...this.allResults];

            if (this.trackerFilter) {
                filtered = filtered.filter(r => r.tracker === this.trackerFilter);
            }

            if (this.sizeFilter && SIZE_LIMITS[this.sizeFilter]) {
                const limits = SIZE_LIMITS[this.sizeFilter];
                filtered = filtered.filter(r => {
                    const size = r.size || 0;
                    if (limits.min && limits.max) return size >= limits.min && size <= limits.max;
                    if (limits.min) return size >= limits.min;
                    if (limits.max) return size < limits.max;
                    return true;
                });
            }

            if (this.qualityFilter) {
                filtered = filtered.filter(r => this.detectQuality(r.title) === this.qualityFilter);
            }

            if (this.searchQuery) {
                const query = this.searchQuery.toLowerCase();
                filtered = filtered.filter(r => r.title.toLowerCase().includes(query));
            }

            this.filteredResults = filtered;
            this.filteredCount = filtered.length;
            this.applySort();
        },

        applySort() {
            const [field, direction] = this.sortBy.split('-');
            this.filteredResults.sort((a, b) => {
                let valA, valB;
                switch (field) {
                    case 'seeders':
                        valA = a.seeders || 0;
                        valB = b.seeders || 0;
                        break;
                    case 'date':
                        valA = new Date(a.publish_date).getTime();
                        valB = new Date(b.publish_date).getTime();
                        break;
                    default:
                        return 0;
                }
                return direction === 'asc' ? valA - valB : valB - valA;
            });
        },

        clearFilters() {
            this.sortBy = 'seeders-desc';
            this.trackerFilter = '';
            this.sizeFilter = '';
            this.qualityFilter = '';
            this.searchQuery = '';
        },

        hasActiveFilters() {
            return this.trackerFilter || this.sizeFilter || this.qualityFilter || this.searchQuery;
        },

        sizeLabel(filter) {
            const labels = {
                small: '< 500 MB',
                medium: '500 MB - 2 GB',
                large: '2 GB - 5 GB',
                huge: '> 5 GB'
            };
            return labels[filter] || filter;
        },

        qualityLabel(filter) {
            const labels = {
                '4k': '4K / UHD', '1080p': '1080p', '720p': '720p',
                'web': 'WEB', 'hdtv': 'HDTV', 'dvd': 'DVD', 'other': 'Other'
            };
            return labels[filter] || filter;
        },

        formatSize(bytes) {
            if (!bytes || bytes === 0) return '0 B';
            const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
            const i = Math.floor(Math.log(bytes) / Math.log(1024));
            return parseFloat((bytes / Math.pow(1024, i)).toFixed(2)) + ' ' + sizes[i];
        },

        formatDate(dateStr) {
            if (!dateStr) return 'Unknown';
            const date = new Date(dateStr);
            return date.toLocaleDateString('en-US', {
                year: 'numeric',
                month: 'short',
                day: 'numeric'
            });
        }
    };
}

package myanimelist

// Custom field, applicable to both Anime.Details() and Manga.Details().
// Transforms into all available fields before sending request.
const FieldAllAvailable string = "*"

// Shared fields for both Anime.Details() and Manga.Details().
const (
	FieldID                string = "id"
	FieldTitle             string = "title"
	FieldMainPicture       string = "main_picture"
	FieldAlternativeTitles string = "alternative_titles"
	FieldStartDate         string = "start_date"
	FieldEndDate           string = "end_date"
	FieldSynopsis          string = "synopsis"
	FieldMean              string = "mean"
	FieldRank              string = "rank"
	FieldPopularity        string = "popularity"
	FieldNumListUsers      string = "num_list_users"
	FieldNumScoringUsers   string = "num_scoring_users"
	FieldNSFW              string = "nsfw"
	FieldCreatedAt         string = "created_at"
	FieldUpdatedAt         string = "updated_at"
	FieldMediaType         string = "media_type"
	FieldStatus            string = "status"
	FieldGenres            string = "genres"
	FieldMyListStatus      string = "my_list_status"
	FieldPictures          string = "pictures"
	FieldBackground        string = "background"
	FieldRelatedAnime      string = "related_anime"
	FieldRelatedManga      string = "related_manga"
	FieldRecommendations   string = "recommendations"
	FieldStudios           string = "studios"
	FieldStatistics        string = "statistics"
)

var generalFields = []string{FieldID, FieldTitle, FieldMainPicture, FieldAlternativeTitles,
	FieldStartDate, FieldEndDate, FieldSynopsis, FieldMean, FieldRank, FieldPopularity,
	FieldNumListUsers, FieldNumScoringUsers, FieldNSFW, FieldCreatedAt, FieldUpdatedAt,
	FieldMediaType, FieldStatus, FieldGenres, FieldMyListStatus, FieldPictures,
	FieldBackground, FieldRelatedAnime, FieldRelatedManga, FieldRecommendations, FieldStatistics}

// Anime.Details() only fields
const (
	FieldNumEpisodes            string = "num_episodes"
	FieldStartSeason            string = "start_season"
	FieldBroadcast              string = "broadcast"
	FieldSource                 string = "source"
	FieldAverageEpisodeDuration string = "average_episode_duration"
	FieldRating                 string = "rating"
)

var animeFields = []string{FieldNumEpisodes, FieldStartSeason, FieldBroadcast, FieldSource,
	FieldAverageEpisodeDuration, FieldRating, FieldStudios}

// Manga.Details() only fields.
const (
	FieldNumVolumes    string = "num_volumes"
	FieldNumChapters   string = "num_chapters"
	FieldAuthors       string = "authors{first_name,last_name}"
	FieldSerialization string = "serialization{name}"
)

var mangaFields = []string{FieldNumVolumes, FieldNumChapters, FieldAuthors, FieldSerialization}

// Shared ranks for both Anime.Top() and Manga.Top().
const (
	// Top Anime|Manga Series
	RankAll string = "all"
	// Top Anime|Manga by Popularity
	RankByPopularity string = "bypopularity"
	// Top Anime|Manga by Favorite
	RankFavorite string = "favorite"
)

var generalRankings = []string{RankAll, RankByPopularity, RankFavorite}

// Anime.Top() only
const (
	// Top Airing Anime
	RankAiring string = "airing"
	// Top Upcoming Anime
	RankUpcoming string = "upcoming"
	// Top Anime TV Series
	RankTV string = "tv"
	// Top Anime OVA Series
	RankOVA string = "ova"
	// Top Anime Movies
	RankMovie string = "movie"
	// Top Anime Specials
	RankSpecials string = "special"
)

var animeRankings = []string{RankAiring, RankUpcoming, RankTV, RankOVA, RankMovie, RankSpecials}

// Manga.Top() only
const (
	// Top Manga
	RankManga string = "manga"
	// Top Novels
	RankNovels string = "novels"
	// Top One-shots
	RankOneShots string = "oneshots"
	// Top Doujinshi
	RankDoujinshi string = "doujin"
	// Top Manhwa
	RankManhwa string = "manhwa"
	// Top Manhua
	RankManhua string = "manhua"
)

var mangaRankings = []string{RankManga, RankNovels, RankOneShots, RankDoujinshi, RankManhwa, RankManhua}

// Shared statuses, working for Anime.List.User() && Manga.List.User()
// and for AnimeConfig.SetStatus() && MangaConfig.SetStatus() too.
const (
	StatusOnHold    string = "on_hold"
	StatusDropped   string = "dropped"
	StatusCompleted string = "completed"
)

var generalStatuses = []string{StatusOnHold, StatusDropped, StatusCompleted}

// Anime-only statuses (Anime.List.User() and AnimeConfig.SetStatus())
const (
	StatusWatching    string = "watching"
	StatusPlanToWatch string = "plan_to_watch"
)

var animeStatuses = []string{StatusWatching, StatusPlanToWatch}

// Manga-only statuses (Manga.List.User() and MangaConfig.SetStatus())
const (
	StatusReading    string = "reading"
	StatusPlanToRead string = "plan_to_read"
)

var mangaStatuses = []string{StatusReading, StatusPlanToRead}

// Predefined season values. Used for 'Anime.Seasonal()'
const (
	// January, February, March
	SeasonWinter string = "winter"
	// April, May, June
	SeasonSpring string = "spring"
	// July, August, September
	SeasonSummer string = "summer"
	// October, November, December
	SeasonFall string = "fall"
)

var seasons = []string{SeasonWinter, SeasonSpring, SeasonSummer, SeasonFall}

// Used to sort Anime.Seasonal() response
const (
	SortByScore      string = "anime_score"
	SortByUsersLists string = "anime_num_list_users"
)

// Anime or manga priority. Used with AnimeConfig.SetPriority() and MangaConfig.SetPriority()
const (
	PriorityLow = iota
	PriorityMedium
	PriorityHigh
)

var priorities = []int{PriorityLow, PriorityMedium, PriorityHigh}

// User list's sort
const (
	SortListByScore      string = "list_score"
	SortListByUpdateDate string = "list_updated_at"
	SortListByTitle      string = "title"
	SortListByStartDate  string = "start_date"
	SortListByID         string = "id"
)

var listSortings = []string{SortListByScore, SortListByUpdateDate, SortListByTitle,
	SortListByStartDate, SortListByID}

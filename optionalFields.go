package myanimelist

// Custom field. Transforms into all available fields before sending request
const FieldAllAvailable string = "*"

// Shared fields
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

var GeneralFields = []string{FieldID, FieldTitle, FieldMainPicture, FieldAlternativeTitles,
	FieldStartDate, FieldEndDate, FieldSynopsis, FieldMean, FieldRank, FieldPopularity,
	FieldNumListUsers, FieldNumScoringUsers, FieldNSFW, FieldCreatedAt, FieldUpdatedAt,
	FieldMediaType, FieldStatus, FieldGenres, FieldMyListStatus, FieldPictures,
	FieldBackground, FieldRelatedAnime, FieldRelatedManga, FieldRecommendations, FieldStatistics}

// Anime-only fields
const (
	FieldNumEpisodes            string = "num_episodes"
	FieldStartSeason            string = "start_season"
	FieldBroadcast              string = "broadcast"
	FieldSource                 string = "source"
	FieldAverageEpisodeDuration string = "average_episode_duration"
	FieldRating                 string = "rating"
)

var AnimeFields = []string{FieldNumEpisodes, FieldStartSeason, FieldBroadcast, FieldSource,
	FieldAverageEpisodeDuration, FieldRating, FieldStudios}

// Manga-only fields
const (
	FieldNumVolumes    string = "num_volumes"
	FieldNumChapters   string = "num_chapters"
	FieldAuthors       string = "authors{first_name,last_name}"
	FieldSerialization string = "serialization{name}"
)

var MangaFields = []string{FieldNumVolumes, FieldNumChapters, FieldAuthors, FieldSerialization}

// Shared ranks
const (
	// Top Anime|Manga Series
	RankAll string = "all"
	// Top Anime|Manga by Popularity
	RankByPopularity string = "bypopularity"
	// Top Anime|Manga by Favorite
	RankFavorite string = "favorite"
)

var GeneralRankings = []string{RankAll, RankByPopularity, RankFavorite}

// Anime related ranks
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

var AnimeRankings = []string{RankAiring, RankUpcoming, RankTV, RankOVA, RankMovie, RankSpecials}

// Manga related ranks
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

var MangaRankings = []string{RankManga, RankNovels, RankOneShots, RankDoujinshi, RankManhwa, RankManhua}

// Shared statuses
const (
	StatusOnHold    string = "on_hold"
	StatusDropped   string = "dropped"
	StatusCompleted string = "completed"
)

var GeneralStatuses = []string{StatusOnHold, StatusDropped, StatusCompleted}

// Anime-only statuses
const (
	StatusWatching    string = "watching"
	StatusPlanToWatch string = "plan_to_watch"
)

var AnimeStatuses = []string{StatusWatching, StatusPlanToWatch}

// Manga-only statuses
const (
	StatusReading    string = "reading"
	StatusPlanToRead string = "plan_to_read"
)

var MangaStatuses = []string{StatusReading, StatusPlanToRead}

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

var Seasons = []string{SeasonWinter, SeasonSpring, SeasonSummer, SeasonFall}

const (
	SortByScore      string = "anime_score"
	SortByUsersLists string = "anime_num_list_users"
)

const (
	PriorityLow = iota
	PriorityMedium
	PriorityHigh
)

var Priorities = []int{PriorityLow, PriorityMedium, PriorityHigh}

const (
	SortListByScore      string = "list_score"
	SortListByUpdateDate string = "list_updated_at"
	SortListByTitle      string = "title"
	SortListByStartDate  string = "start_date"
	SortListByID         string = "id"
)

var ListSortings = []string{SortListByScore, SortListByUpdateDate, SortListByTitle,
	SortListByStartDate, SortListByID}

package usecase

import (
	"repertoire/server/domain/usecase/album"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/domain/usecase/artist/band/member"
	"repertoire/server/domain/usecase/playlist"
	playlistSong "repertoire/server/domain/usecase/playlist/song"
	"repertoire/server/domain/usecase/search"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/domain/usecase/udata/band/member/role"
	"repertoire/server/domain/usecase/udata/guitar/tuning"
	"repertoire/server/domain/usecase/udata/instrument"
	"repertoire/server/domain/usecase/udata/section/types"
	"repertoire/server/domain/usecase/user"

	"go.uber.org/fx"
)

var albumUseCases = fx.Options(
	fx.Provide(album.NewAddSongsToAlbum),
	fx.Provide(album.NewAddPerfectRehearsalsToAlbums),
	fx.Provide(album.NewBulkDeleteAlbums),
	fx.Provide(album.NewCreateAlbum),
	fx.Provide(album.NewDeleteAlbum),
	fx.Provide(album.NewDeleteImageFromAlbum),
	fx.Provide(album.NewGetAlbum),
	fx.Provide(album.NewGetAlbumFiltersMetadata),
	fx.Provide(album.NewGetAllAlbums),
	fx.Provide(album.NewMoveSongFromAlbum),
	fx.Provide(album.NewRemoveSongsFromAlbum),
	fx.Provide(album.NewSaveImageToAlbum),
	fx.Provide(album.NewUpdateAlbum),
)

var artistUseCases = fx.Options(
	fx.Provide(artist.NewAddAlbumsToArtist),
	fx.Provide(artist.NewAddPerfectRehearsalsToArtists),
	fx.Provide(artist.NewAddSongsToArtist),
	fx.Provide(artist.NewBulkDeleteArtists),
	fx.Provide(artist.NewCreateArtist),
	fx.Provide(artist.NewDeleteArtist),
	fx.Provide(artist.NewDeleteImageFromArtist),
	fx.Provide(artist.NewGetAllArtists),
	fx.Provide(artist.NewGetArtist),
	fx.Provide(artist.NewGetArtistFiltersMetadata),
	fx.Provide(artist.NewRemoveAlbumsFromArtist),
	fx.Provide(artist.NewRemoveSongsFromArtist),
	fx.Provide(artist.NewSaveImageToArtist),
	fx.Provide(artist.NewUpdateArtist),

	fx.Provide(member.NewCreateBandMember),
	fx.Provide(member.NewDeleteBandMember),
	fx.Provide(member.NewDeleteImageFromBandMember),
	fx.Provide(member.NewMoveBandMember),
	fx.Provide(member.NewSaveImageToBandMember),
	fx.Provide(member.NewUpdateBandMember),

	fx.Provide(member.NewGetBandMemberRoles),
)

var playlistUseCases = fx.Options(
	fx.Provide(playlist.NewAddAlbumsToPlaylist),
	fx.Provide(playlist.NewAddArtistsToPlaylist),
	fx.Provide(playlist.NewBulkDeletePlaylists),
	fx.Provide(playlist.NewCreatePlaylist),
	fx.Provide(playlist.NewDeletePlaylist),
	fx.Provide(playlist.NewDeleteImageFromPlaylist),
	fx.Provide(playlist.NewGetAllPlaylists),
	fx.Provide(playlist.NewGetPlaylist),
	fx.Provide(playlist.NewGetPlaylistFiltersMetadata),
	fx.Provide(playlist.NewSaveImageToPlaylist),
	fx.Provide(playlist.NewUpdatePlaylist),

	fx.Provide(playlistSong.NewAddSongsToPlaylist),
	fx.Provide(playlistSong.NewGetPlaylistSongs),
	fx.Provide(playlistSong.NewMoveSongFromPlaylist),
	fx.Provide(playlistSong.NewRemoveSongsFromPlaylist),
	fx.Provide(playlistSong.NewShufflePlaylistSongs),
)

var searchUseCases = fx.Options(
	fx.Provide(search.NewGet),
	fx.Provide(search.NewMeiliWebhook),
)

var songUseCases = fx.Options(
	fx.Provide(song.NewAddPerfectSongRehearsal),
	fx.Provide(song.NewAddPerfectSongRehearsals),
	fx.Provide(song.NewAddPartialSongRehearsal),
	fx.Provide(song.NewBulkDeleteSongs),
	fx.Provide(song.NewCreateSong),
	fx.Provide(song.NewDeleteSong),
	fx.Provide(song.NewDeleteImageFromSong),
	fx.Provide(song.NewGetAllSongs),
	fx.Provide(song.NewGetSong),
	fx.Provide(song.NewGetSongFiltersMetadata),
	fx.Provide(song.NewSaveImageToSong),
	fx.Provide(song.NewUpdateSong),
	fx.Provide(song.NewUpdateSongSettings),

	fx.Provide(song.NewGetGuitarTunings),
	fx.Provide(song.NewGetInstruments),
	fx.Provide(section.NewGetSongSectionTypes),

	fx.Provide(section.NewCreateSongSection),
	fx.Provide(section.NewDeleteSongSection),
	fx.Provide(section.NewMoveSongSection),
	fx.Provide(section.NewUpdateAllSongSections),
	fx.Provide(section.NewUpdateSongSection),
	fx.Provide(section.NewUpdateSongSectionsOccurrences),
	fx.Provide(section.NewUpdateSongSectionsPartialOccurrences),
)

var userDataUseCases = fx.Options(
	fx.Provide(role.NewCreateBandMemberRole),
	fx.Provide(role.NewDeleteBandMemberRole),
	fx.Provide(role.NewMoveBandMemberRole),

	fx.Provide(tuning.NewCreateGuitarTuning),
	fx.Provide(tuning.NewDeleteGuitarTuning),
	fx.Provide(tuning.NewMoveGuitarTuning),

	fx.Provide(instrument.NewCreateInstrument),
	fx.Provide(instrument.NewDeleteInstrument),
	fx.Provide(instrument.NewMoveInstrument),

	fx.Provide(types.NewCreateSongSectionType),
	fx.Provide(types.NewDeleteSongSectionType),
	fx.Provide(types.NewMoveSongSectionType),
)

var userUseCases = fx.Options(
	fx.Provide(user.NewDeleteUser),
	fx.Provide(user.NewDeleteProfilePictureFromUser),
	fx.Provide(user.NewGetUser),
	fx.Provide(user.NewSaveProfilePictureToUser),
	fx.Provide(user.NewSignUp),
	fx.Provide(user.NewUpdateUser),
)

var Module = fx.Options(
	albumUseCases,
	artistUseCases,
	playlistUseCases,
	searchUseCases,
	songUseCases,
	userDataUseCases,
	userUseCases,
)

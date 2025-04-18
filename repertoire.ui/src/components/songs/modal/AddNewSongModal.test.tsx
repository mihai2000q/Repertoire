import { reduxRender, withToastify } from '../../../test-utils.tsx'
import AddNewSongModal from './AddNewSongModal.tsx'
import { vi } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { act, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import { CreateSongRequest } from '../../../types/requests/SongRequests.ts'
import { GuitarTuning, SongSectionType } from '../../../types/models/Song.ts'
import Difficulty from '../../../types/enums/Difficulty.ts'
import { AlbumSearch, ArtistSearch } from '../../../types/models/Search.ts'
import SearchType from '../../../types/enums/SearchType.ts'
import dayjs from 'dayjs'

describe('Add New Song Modal', () => {
  const albums: AlbumSearch[] = [
    {
      id: '1',
      title: 'Album 1',
      imageUrl: 'something-album.png',
      artist: {
        id: '1',
        name: 'Some Artist'
      },
      releaseDate: '2024-10-12T10:30',
      type: SearchType.Album
    },
    {
      id: '2',
      title: 'Album 2',
      type: SearchType.Album
    },
    {
      id: '3',
      title: 'Album 3',
      type: SearchType.Album
    }
  ]

  const artists: ArtistSearch[] = [
    {
      id: '1',
      name: 'Artist 1',
      type: SearchType.Artist
    },
    {
      id: '2',
      name: 'Artist 2',
      type: SearchType.Artist
    },
    {
      id: '3',
      name: 'Artist 3',
      type: SearchType.Artist
    }
  ]

  const guitarTunings: GuitarTuning[] = [
    {
      id: '1',
      name: 'E Standard'
    },
    {
      id: '2',
      name: 'Drop D'
    },
    {
      id: '3',
      name: 'Drop A'
    }
  ]

  const songSectionTypes: SongSectionType[] = [
    {
      id: '1',
      name: 'Solo'
    },
    {
      id: '2',
      name: 'Riff'
    }
  ]

  const validSongsterrLink =
    'https://www.songsterr.com/a/wsa/metallica-master-of-puppets-tab-s455118'
  const validYoutubeLink = 'https://www.youtube.com/watch?v=E0ozmU9cJDg'

  const handlers = [
    http.get('/search', async (req) => {
      const type = new URL(req.request.url).searchParams.get('type')
      const albumResponse: WithTotalCountResponse<AlbumSearch> = {
        models: albums,
        totalCount: albums.length
      }
      const artistResponse: WithTotalCountResponse<ArtistSearch> = {
        models: artists,
        totalCount: artists.length
      }
      return HttpResponse.json(type === SearchType.Album ? albumResponse : artistResponse)
    }),
    http.get('/songs/guitar-tunings', async () => {
      return HttpResponse.json(guitarTunings)
    }),
    http.get('/songs/sections/types', async () => {
      return HttpResponse.json(songSectionTypes)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  describe('should render', () => {
    function shouldRender() {
      expect(screen.getByRole('dialog', { name: /add new song/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /add new song/i })).toBeInTheDocument()

      expect(screen.getByRole('button', { name: /first step/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /second step/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /final step/i })).toBeInTheDocument()
    }

    it('first step', () => {
      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      shouldRender()

      expect(screen.getByRole('textbox', { name: /title/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /title/i })).not.toBeInvalid()

      expect(screen.getByRole('textbox', { name: /album/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /album/i })).toHaveValue('')

      expect(screen.getByRole('textbox', { name: /artist/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /artist/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /artist/i })).not.toBeDisabled()

      expect(screen.getByRole('textbox', { name: /description/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /description/i })).toHaveValue('')

      expect(screen.getByRole('button', { name: /next/i })).toBeInTheDocument()
    })

    it('second step', async () => {
      const user = userEvent.setup()

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      await user.type(screen.getByRole('textbox', { name: /title/i }), 'something') // to validate and go to next step
      await user.click(screen.getByRole('button', { name: /second step/i }))

      shouldRender()

      expect(screen.getByRole('textbox', { name: /guitar tuning/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /guitar tuning/i })).toHaveValue('')

      expect(screen.getByRole('textbox', { name: /difficulty/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /difficulty/i })).toHaveValue('')

      expect(screen.getByRole('textbox', { name: /bpm/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /bpm/i })).toHaveValue('')

      expect(screen.getByRole('button', { name: /release date/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /release date/i })).toHaveTextContent(
        new RegExp('release date', 'i')
      )

      expect(screen.getByRole('button', { name: /add section/i })).toBeInTheDocument()

      expect(screen.getByRole('button', { name: /previous/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /next/i })).toBeInTheDocument()
    })

    it('final step', async () => {
      const user = userEvent.setup()

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      await user.type(screen.getByRole('textbox', { name: /title/i }), 'something') // to validate and go to next step
      await user.click(screen.getByRole('button', { name: /final step/i }))

      shouldRender()

      expect(screen.getByRole('textbox', { name: /youtube/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /youtube/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /youtube/i })).not.toBeInvalid()

      expect(screen.getByRole('textbox', { name: /songsterr/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /songsterr/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /songsterr/i })).not.toBeInvalid()

      expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()

      expect(screen.getByRole('button', { name: /previous/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /submit/i })).toBeInTheDocument()
    })
  })

  describe('first step', () => {
    it('should validate title on blur', async () => {
      const user = userEvent.setup()

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      const titleField = screen.getByRole('textbox', { name: /title/i })

      await user.type(titleField, 'something')
      expect(titleField).not.toBeInvalid()

      await user.clear(titleField)
      act(() => titleField.blur())
      expect(titleField).toBeInvalid()
    })

    it('should validate on clicking steps', async () => {
      const user = userEvent.setup()

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      const titleField = screen.getByRole('textbox', { name: /title/i })

      await user.click(screen.getByRole('button', { name: /second step/i }))
      expect(titleField).toBeInvalid()

      await user.type(titleField, 'something')
      expect(titleField).not.toBeInvalid()

      await user.click(screen.getByRole('button', { name: /final step/i }))
      expect(titleField).not.toBeInTheDocument()
    })

    it('should validate on clicking next', async () => {
      const user = userEvent.setup()

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      const titleField = screen.getByRole('textbox', { name: /title/i })

      await user.click(screen.getByRole('button', { name: /next/i }))
      expect(titleField).toBeInvalid()

      await user.type(titleField, 'something')
      expect(titleField).not.toBeInvalid()

      await user.click(screen.getByRole('button', { name: /next/i }))
      expect(titleField).not.toBeInTheDocument()
    })
  })

  describe('second step', () => {
    async function goToSecondStep() {
      const user = userEvent.setup()

      await user.type(screen.getByRole('textbox', { name: /title/i }), 'something')
      await user.click(screen.getByRole('button', { name: /next/i }))
    }

    it('should add new section fields by clicking on add section and remove each one by clicking on remove', async () => {
      const user = userEvent.setup()

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      await goToSecondStep()

      expect(screen.queryByRole('button', { name: 'drag-handle' })).not.toBeInTheDocument()
      expect(screen.queryByRole('textbox', { name: /type/i })).not.toBeInTheDocument()
      expect(screen.queryByRole('textbox', { name: /name/i })).not.toBeInTheDocument()

      // add one section
      await user.click(screen.getByRole('button', { name: /add section/i }))

      expect(screen.getByRole('button', { name: 'drag-handle' })).toBeInTheDocument()

      expect(screen.getByRole('textbox', { name: /type/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /type/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /type/i })).not.toBeInvalid()

      expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()

      expect(screen.getByRole('button', { name: 'remove-section' })).toBeInTheDocument()

      // add second section
      await user.click(screen.getByRole('button', { name: /add section/i }))
      expect(screen.queryAllByRole('textbox', { name: /name/i })).toHaveLength(2)

      // populate first section's name
      await user.type(screen.queryAllByRole('textbox', { name: /name/i })[0], 'something')

      // remove first section
      await user.click(screen.queryAllByRole('button', { name: 'remove-section' })[0])
      expect(screen.queryAllByRole('textbox', { name: /name/i })).toHaveLength(1)
      expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('') // not the one from above
    })

    it("should validate sections' name and type on blur", async () => {
      const user = userEvent.setup()

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      await goToSecondStep()

      await user.click(screen.getByRole('button', { name: /add section/i }))

      const nameField = screen.getByRole('textbox', { name: /name/i })

      await user.type(nameField, 'some name')
      expect(nameField).not.toBeInvalid()

      await user.clear(nameField)
      act(() => nameField.blur())
      expect(nameField).toBeInvalid()
    })

    it('should validate on clicking steps', async () => {
      const user = userEvent.setup()

      const newName = 'some name'
      const newType = songSectionTypes[0]

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      await goToSecondStep()

      await user.click(screen.getByRole('button', { name: /add section/i }))

      const nameField = screen.getByRole('textbox', { name: /name/i })
      const typeField = screen.getByRole('textbox', { name: /type/i })

      await user.click(screen.getByRole('button', { name: /first step/i }))
      expect(nameField).toBeInvalid()
      expect(typeField).toBeInvalid()

      await user.type(nameField, newName)
      expect(nameField).not.toBeInvalid()

      await user.click(typeField)
      await user.click(await screen.findByRole('option', { name: newType.name }))
      expect(typeField).not.toBeInvalid()

      await user.click(screen.getByRole('button', { name: /final step/i }))
      expect(nameField).not.toBeInTheDocument()
      expect(typeField).not.toBeInTheDocument()
    })

    it('should validate on clicking next', async () => {
      const user = userEvent.setup()

      const newName = 'some name'
      const newType = songSectionTypes[0]

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      await goToSecondStep()

      await user.click(screen.getByRole('button', { name: /add section/i }))

      const nameField = screen.getByRole('textbox', { name: /name/i })
      const typeField = screen.getByRole('textbox', { name: /type/i })

      await user.click(screen.getByRole('button', { name: /previous/i }))
      expect(nameField).toBeInvalid()
      expect(typeField).toBeInvalid()
      await user.click(screen.getByRole('button', { name: /next/i }))
      expect(nameField).toBeInvalid()
      expect(typeField).toBeInvalid()

      await user.type(nameField, newName)
      expect(nameField).not.toBeInvalid()

      await user.click(typeField)
      await user.click(await screen.findByRole('option', { name: newType.name }))
      expect(typeField).not.toBeInvalid()

      await user.click(screen.getByRole('button', { name: /next/i }))
      expect(nameField).not.toBeInTheDocument()
      expect(typeField).not.toBeInTheDocument()
    })

    it.skip('should be able to reorder sections', () => {})
  })

  describe('final step', () => {
    async function goToFinalStep() {
      const user = userEvent.setup()

      await user.type(screen.getByRole('textbox', { name: /title/i }), 'something')
      await user.click(screen.getByRole('button', { name: /final step/i }))
    }

    it('should validate songsterr and youtube on blur', async () => {
      const user = userEvent.setup()

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      await goToFinalStep()

      const songsterrField = screen.getByRole('textbox', { name: /songsterr/i })
      const youtubeField = screen.getByRole('textbox', { name: /youtube/i })

      // songsterr
      await user.type(songsterrField, 'something')
      act(() => songsterrField.blur())
      expect(songsterrField).toBeInvalid()

      await user.click(songsterrField)
      await user.paste(validSongsterrLink)
      expect(songsterrField).not.toBeInvalid()

      await user.clear(songsterrField)
      expect(songsterrField).not.toBeInvalid()

      // youtube
      await user.type(youtubeField, 'something')
      act(() => youtubeField.blur())
      expect(youtubeField).toBeInvalid()

      await user.click(youtubeField)
      await user.paste(validYoutubeLink)
      expect(youtubeField).not.toBeInvalid()

      await user.clear(youtubeField)
      expect(youtubeField).not.toBeInvalid()
    })

    it('should validate on clicking steps', async () => {
      const user = userEvent.setup()

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      await goToFinalStep()

      const songsterrField = screen.getByRole('textbox', { name: /songsterr/i })

      await user.type(songsterrField, 'something')
      await user.click(screen.getByRole('button', { name: /first step/i }))
      expect(songsterrField).toBeInvalid()

      await user.clear(songsterrField)
      await user.click(songsterrField)
      await user.paste(validSongsterrLink)
      expect(songsterrField).not.toBeInvalid()

      await user.click(screen.getByRole('button', { name: /first step/i }))
      expect(songsterrField).not.toBeInTheDocument()
    })

    it('should validate on clicking previous or submit', async () => {
      const user = userEvent.setup()

      reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

      await goToFinalStep()

      const songsterrField = screen.getByRole('textbox', { name: /songsterr/i })

      // check previous button validation
      await user.type(songsterrField, 'something')
      await user.click(screen.getByRole('button', { name: /previous/i }))
      expect(songsterrField).toBeInvalid()

      // check submit button validation
      await user.clear(songsterrField)
      expect(songsterrField).not.toBeInvalid()
      await user.type(songsterrField, 'something')
      await user.click(screen.getByRole('button', { name: /submit/i }))
      expect(songsterrField).toBeInvalid()

      // check previous button validation
      await user.click(songsterrField)
      await user.paste(validSongsterrLink)
      expect(songsterrField).not.toBeInvalid()

      await user.click(screen.getByRole('button', { name: /previous/i }))
      expect(songsterrField).not.toBeInTheDocument()
    })
  })

  it('should be able to change steps', async () => {
    const user = userEvent.setup()

    reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

    await user.type(screen.getByRole('textbox', { name: /title/i }), 'something') // so that validation passes

    function expectFirstStep() {
      expect(screen.getByRole('textbox', { name: /title/i })).toBeInTheDocument()
      expect(screen.queryByRole('button', { name: /previous/i })).not.toBeInTheDocument()
      expect(screen.getByRole('button', { name: /next/i })).toBeInTheDocument()
    }

    function expectSecondStep() {
      expect(screen.queryByRole('textbox', { name: /title/i })).not.toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /guitar tuning/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /previous/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /next/i })).toBeInTheDocument()
    }

    function expectFinalStep() {
      expect(screen.queryByRole('textbox', { name: /title/i })).not.toBeInTheDocument()
      expect(screen.queryByRole('textbox', { name: /guitar tuning/i })).not.toBeInTheDocument()
      expect(screen.getByRole('textbox', { name: /songsterr/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /previous/i })).toBeInTheDocument()
      expect(screen.queryByRole('button', { name: /next/i })).not.toBeInTheDocument()
      expect(screen.getByRole('button', { name: /submit/i })).toBeInTheDocument()
    }

    // use next and previous buttons
    expectFirstStep()

    await user.click(screen.getByRole('button', { name: /next/i }))
    expectSecondStep()

    await user.click(screen.getByRole('button', { name: /previous/i }))
    expectFirstStep()

    await user.click(screen.getByRole('button', { name: /next/i }))
    await user.click(screen.getByRole('button', { name: /next/i }))
    expectFinalStep()

    await user.click(screen.getByRole('button', { name: /previous/i }))
    expectSecondStep()

    await user.click(screen.getByRole('button', { name: /previous/i }))
    expectFirstStep()

    // use step buttons
    await user.click(screen.getByRole('button', { name: /second step/i }))
    expectSecondStep()

    await user.click(screen.getByRole('button', { name: /final step/i }))
    expectFinalStep()

    await user.click(screen.getByRole('button', { name: /first step/i }))
    expectFirstStep()
  })

  it('should disable the artist when an album is selected', async () => {
    const user = userEvent.setup()

    const newAlbum = albums[1]

    reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

    await user.click(screen.getByRole('textbox', { name: /album/i }))
    await user.click(await screen.findByRole('option', { name: newAlbum.title }))
    expect(screen.getByRole('textbox', { name: /artist/i })).toBeDisabled()
  })

  it('should show artist info icon that it will be inherited from album, when album is selected and it does not have an artist', async () => {
    const user = userEvent.setup()

    const newAlbum = albums[1]

    reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

    await user.click(screen.getByRole('textbox', { name: /album/i }))
    await user.click(await screen.findByRole('option', { name: newAlbum.title }))
    expect(await screen.findByLabelText('artist-info')).toBeInTheDocument()
  })

  it('should show artist, release date and image info icons that they will be inherited from album, when the selected album has those values', async () => {
    const user = userEvent.setup()

    const newAlbum = albums[0]

    reduxRender(<AddNewSongModal opened={true} onClose={() => {}} />)

    await user.type(screen.getByRole('textbox', { name: /title/i }), 'something') // so that validation passes

    await user.click(screen.getByRole('textbox', { name: /album/i }))
    await user.click(await screen.findByRole('option', { name: newAlbum.title }))
    expect(await screen.findByLabelText('artist-info')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: /second step/i }))
    expect(screen.getByLabelText('release-date-info')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: /final step/i }))
    expect(screen.getByText(/if no image is uploaded, it will be inherited/i)).toBeInTheDocument()
  })

  describe('should send create requests', () => {
    async function expectFormIsReset() {
      const user = userEvent.setup()

      // first step
      expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /album/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /artist/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /description/i })).toHaveValue('')

      await user.type(screen.getByRole('textbox', { name: /title/i }), 'something') // pass validation

      // second step
      await user.click(screen.getByRole('button', { name: /second step/i }))
      expect(screen.getByRole('textbox', { name: /guitar tuning/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /difficulty/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /bpm/i })).toHaveValue('')
      expect(screen.getByRole('button', { name: /release date/i })).toHaveTextContent(
        new RegExp('release date', 'i')
      )
      expect(screen.queryAllByRole('textbox', { name: /name/i })).toHaveLength(0)

      // final step
      await user.click(screen.getByRole('button', { name: /final step/i }))
      expect(screen.getByRole('textbox', { name: /songsterr/i })).toHaveValue('')
      expect(screen.getByRole('textbox', { name: /youtube/i })).toHaveValue('')
      expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
    }

    it('minimal - only title is filled', async () => {
      const user = userEvent.setup()

      const newTitle = 'New Title'

      const onClose = vi.fn()

      let capturedRequest: CreateSongRequest
      server.use(
        http.post('/songs', async (req) => {
          capturedRequest = (await req.request.json()) as CreateSongRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(withToastify(<AddNewSongModal opened={true} onClose={onClose} />))

      await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
      await user.click(screen.getByRole('button', { name: /final step/i }))
      await user.click(screen.getByRole('button', { name: /submit/i }))

      expect(capturedRequest).toStrictEqual({
        title: newTitle,
        description: '',
        sections: []
      })
      expect(onClose).toHaveBeenCalledOnce()

      expect(screen.getByText(new RegExp(`${newTitle} added`, 'i'))).toBeInTheDocument()
      await expectFormIsReset()
    })

    it('with album selected', async () => {
      const user = userEvent.setup()

      const newTitle = 'New Title'
      const newAlbum = albums[0]

      const onClose = vi.fn()

      let capturedRequest: CreateSongRequest
      server.use(
        http.post('/songs', async (req) => {
          capturedRequest = (await req.request.json()) as CreateSongRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(withToastify(<AddNewSongModal opened={true} onClose={onClose} />))

      await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)

      await user.click(screen.getByRole('textbox', { name: /album/i }))
      await user.click(await screen.findByRole('option', { name: newAlbum.title }))

      await user.click(screen.getByRole('button', { name: /final step/i }))
      await user.click(screen.getByRole('button', { name: /submit/i }))

      expect(capturedRequest).toEqual({
        title: newTitle,
        description: '',
        albumId: newAlbum.id,
        sections: []
      })
      expect(onClose).toHaveBeenCalledOnce()

      expect(screen.getByText(new RegExp(`${newTitle} added`, 'i'))).toBeInTheDocument()
      await expectFormIsReset()
    })

    it('with artist selected and new album', async () => {
      const user = userEvent.setup()

      const newTitle = 'New Title'
      const newAlbumTitle = 'New Album Title'
      const newArtist = artists[0]

      const onClose = vi.fn()

      let capturedRequest: CreateSongRequest
      server.use(
        http.post('/songs', async (req) => {
          capturedRequest = (await req.request.json()) as CreateSongRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(withToastify(<AddNewSongModal opened={true} onClose={onClose} />))

      await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
      await user.type(screen.getByRole('textbox', { name: /album/i }), newAlbumTitle)

      await user.click(screen.getByRole('textbox', { name: /artist/i }))
      await user.click(await screen.findByRole('option', { name: newArtist.name }))

      await user.click(screen.getByRole('button', { name: /final step/i }))
      await user.click(screen.getByRole('button', { name: /submit/i }))

      expect(capturedRequest).toEqual({
        title: newTitle,
        description: '',
        artistId: newArtist.id,
        albumTitle: newAlbumTitle,
        sections: []
      })
      expect(onClose).toHaveBeenCalledOnce()

      expect(screen.getByText(new RegExp(`${newTitle} added`, 'i'))).toBeInTheDocument()
      await expectFormIsReset()
    })

    it('maximal with new artist and new album', async () => {
      const user = userEvent.setup()

      // first step data
      const newTitle = 'New Title'
      const newAlbumTitle = 'New Album Title'
      const newArtistName = 'New Artist Name'
      const newDescription = 'New Description with very much significance'

      // second step data
      const newGuitarTuning = guitarTunings[0]
      const newDifficulty = Difficulty.Easy
      const newBpm = 123
      const now = new Date()
      const newReleaseDate = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0, 0)

      // sections
      const newSectionType = songSectionTypes[0]
      const newSectionName = 'New Section'

      const onClose = vi.fn()

      let capturedRequest: CreateSongRequest
      server.use(
        http.post('/songs', async (req) => {
          capturedRequest = (await req.request.json()) as CreateSongRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(withToastify(<AddNewSongModal opened={true} onClose={onClose} />))

      // first step
      await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)
      await user.type(screen.getByRole('textbox', { name: /album/i }), newAlbumTitle)
      await user.type(screen.getByRole('textbox', { name: /artist/i }), newArtistName)
      await user.type(screen.getByRole('textbox', { name: /description/i }), newDescription)

      // second step
      await user.click(screen.getByRole('button', { name: /next/i }))

      await user.click(screen.getByRole('textbox', { name: /guitar tuning/i }))
      await user.click(await screen.findByRole('option', { name: newGuitarTuning.name }))

      await user.click(screen.getByRole('textbox', { name: /difficulty/i }))
      await user.click(await screen.findByRole('option', { name: new RegExp(newDifficulty, 'i') }))

      await user.type(screen.getByRole('textbox', { name: /bpm/i }), newBpm.toString())

      await user.click(screen.getByRole('button', { name: /release date/i }))
      await user.click(
        screen.getByRole('button', { name: dayjs(newReleaseDate).format('D MMMM YYYY') })
      )

      // sections
      await user.click(screen.getByRole('button', { name: /add section/i }))

      await user.click(screen.getByRole('textbox', { name: /type/i }))
      await user.click(await screen.findByRole('option', { name: newSectionType.name }))

      await user.type(screen.getByRole('textbox', { name: /name/i }), newSectionName)

      // final step
      await user.click(screen.getByRole('button', { name: /next/i }))

      await user.click(screen.getByRole('textbox', { name: /songsterr/i }))
      await user.paste(validSongsterrLink)

      await user.click(screen.getByRole('textbox', { name: /youtube/i }))
      await user.paste(validYoutubeLink)

      await user.click(screen.getByRole('button', { name: /submit/i }))

      expect(capturedRequest).toEqual({
        title: newTitle,
        albumTitle: newAlbumTitle,
        artistName: newArtistName,
        description: newDescription,
        guitarTuningId: newGuitarTuning.id,
        difficulty: newDifficulty,
        bpm: newBpm,
        releaseDate: newReleaseDate.toISOString(),
        sections: [
          {
            typeId: newSectionType.id,
            name: newSectionName
          }
        ],
        songsterrLink: validSongsterrLink,
        youtubeLink: validYoutubeLink
      })
      expect(onClose).toHaveBeenCalledOnce()

      expect(screen.getByText(new RegExp(`${newTitle} added`, 'i'))).toBeInTheDocument()
      await expectFormIsReset()
    })
  })

  it('should send create request with save image request, when selecting an image too', async () => {
    const user = userEvent.setup()

    const newTitle = 'New Title'
    const newImage = new File(['something'], 'image.png', { type: 'image/png' })

    const onClose = vi.fn()
    const returnedId = 'the-song-id'

    let capturedRequest: CreateSongRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.post('/songs', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongRequest
        return HttpResponse.json({ id: returnedId })
      }),
      http.put('/songs/images', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(withToastify(<AddNewSongModal opened={true} onClose={onClose} />))

    await user.type(screen.getByRole('textbox', { name: /title/i }), newTitle)

    await user.click(screen.getByRole('button', { name: /final step/i }))
    await user.upload(screen.getByTestId('image-dropzone-input'), newImage)

    await user.click(screen.getByRole('button', { name: /submit/i }))

    expect(capturedRequest).toEqual({
      title: newTitle,
      description: '',
      sections: []
    })
    expect(capturedSaveImageFormData.get('image')).toBeFormDataImage(newImage)
    expect(capturedSaveImageFormData.get('id')).toBe(returnedId)
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(new RegExp(`${newTitle} added`, 'i'))).toBeInTheDocument()

    // form reset
    expect(screen.getByRole('textbox', { name: /title/i })).toHaveValue('')
    await user.type(screen.getByRole('textbox', { name: /title/i }), 'something')
    await user.click(screen.getByRole('button', { name: /final step/i }))
    expect(screen.getByRole('presentation', { name: 'image-dropzone' })).toBeInTheDocument()
  })
})

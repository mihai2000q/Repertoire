import {
  emptySongSection,
  emptySongSettings,
  reduxRender,
  withToastify
} from '../../../../../../../test-utils.tsx'
import AddNewSongPart from './AddNewSongPart.tsx'
import { screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import { Instrument, SongSection } from '../../../../../../../types/models/Song.ts'
import { CreateSongPartRequest } from '../../parts/types/requests/SongPartRequests.ts'
import { BandMember } from '../../../../../../../types/models/Artist.ts'

describe('Add New Song Part', () => {
  const instruments: Instrument[] = [
    {
      id: '1',
      name: 'Guitar'
    },
    {
      id: '2',
      name: 'Piano'
    },
    {
      id: '3',
      name: 'Flute'
    }
  ]

  const bandMembers: BandMember[] = [
    {
      id: '1',
      name: 'Chester',
      roles: []
    },
    {
      id: '2',
      name: 'Michael',
      roles: []
    }
  ]

  const songSections: SongSection[] = [
    {
      ...emptySongSection,
      id: '1',
      name: 'Chorus'
    },
    {
      ...emptySongSection,
      id: '2',
      name: 'Verse'
    }
  ]

  const handlers = [
    http.get('/songs/instruments', async () => {
      return HttpResponse.json(instruments)
    }),
    http.get('/songs/sections', async () => {
      return HttpResponse.json(songSections)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: { songId: '', isArtistBand: false, settings: emptySongSettings }
    })

    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-instrument' })).toBeInTheDocument()
    expect(screen.getByRole('combobox', { name: 'song-section' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(await screen.findByRole('textbox', { name: /name/i })).toHaveFocus()
    expect(screen.getByRole('button', { name: /add/i })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
  })

  it('should disable band member select when the artist is not band', () => {
    reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: {
        songId: '',
        isArtistBand: false,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeDisabled()
  })

  it('should disable band member select when the artist is band, but a song section is not selected', () => {
    reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: {
        songId: '',
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeDisabled()
  })

  it('should enable band member select when the artist is band and a song section is selected', async () => {
    const user = userEvent.setup()

    reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: {
        songId: '',
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    await user.click(screen.getByRole('combobox', { name: 'song-section' }))
    await user.click(await screen.findByRole('option', { name: songSections[0].name }))

    expect(screen.getByRole('button', { name: 'select-band-member' })).not.toBeDisabled()
  })

  it('should have default options based on settings', async () => {
    const user = userEvent.setup()

    const defaultInstrument = instruments[1]
    const defaultBandMember = bandMembers[1]

    const songSettings = { ...emptySongSettings, defaultBandMember, defaultInstrument }

    reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: {
        songId: '',
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: songSettings
      }
    })

    await user.click(screen.getByRole('combobox', { name: 'song-section' }))
    await user.click(await screen.findByRole('option', { name: songSections[0].name }))

    expect(screen.getByRole('button', { name: defaultBandMember.name })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: defaultInstrument.name })).toBeInTheDocument()
  })

  it('should send create request when name is typed', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const songId = 'some id'

    const newName = 'Part 1'

    let capturedRequest: CreateSongPartRequest
    server.use(
      http.post('/songs/parts', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongPartRequest
        return HttpResponse.json({ message: 'part added!' })
      })
    )

    reduxRender(withToastify(<AddNewSongPart opened={true} onClose={onClose} />), {
      song: {
        songId: songId,
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(capturedRequest).toStrictEqual({
      name: newName,
      songId: songId
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newName} added!`)).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
  })

  it('should send create request when all fields are filled', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const songId = 'some id'

    const newName = 'Part 1'
    const newSongSection = songSections[0]
    const newInstrument = instruments[0]
    const newBandMember = bandMembers[0]

    let capturedRequest: CreateSongPartRequest
    server.use(
      http.post('/songs/parts', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongPartRequest
        return HttpResponse.json({ message: 'part added!' })
      })
    )

    reduxRender(withToastify(<AddNewSongPart opened={true} onClose={onClose} />), {
      song: {
        songId: songId,
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    // fill fields
    await user.click(screen.getByRole('combobox', { name: 'song-section' }))
    await user.click(await screen.findByRole('option', { name: newSongSection.name }))

    await user.click(screen.getByRole('button', { name: 'select-band-member' }))
    await user.click(await screen.findByRole('option', { name: newBandMember.name }))

    await user.click(screen.getByRole('button', { name: 'select-instrument' }))
    await user.click(await screen.findByRole('option', { name: newInstrument.name }))

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(capturedRequest).toStrictEqual({
      name: newName,
      songId: songId,
      instrumentId: newInstrument.id,
      sectionId: newSongSection.id,
      bandMemberId: newBandMember.id
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newName} added!`)).toBeInTheDocument()

    // reset fields
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-instrument' })).toBeInTheDocument()
    expect(screen.getByRole('combobox', { name: 'song-section' })).toHaveValue('')
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
  })

  it('should send create request when there are default settings', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const songId = 'some id'
    const settings = {
      ...emptySongSettings,
      defaultBandMember: bandMembers[0],
      defaultInstrument: instruments[0]
    }

    const newName = 'Part 1'
    const newSongSection = songSections[0]
    const newInstrument = instruments[1]

    let capturedRequest: CreateSongPartRequest
    server.use(
      http.post('/songs/parts', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongPartRequest
        return HttpResponse.json({ message: 'part added!' })
      })
    )

    reduxRender(withToastify(<AddNewSongPart opened={true} onClose={onClose} />), {
      song: {
        songId: songId,
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: settings
      }
    })

    // fill fields
    await user.click(screen.getByRole('combobox', { name: 'song-section' }))
    await user.click(await screen.findByRole('option', { name: newSongSection.name }))

    await user.click(screen.getByRole('button', { name: settings.defaultInstrument.name }))
    await user.clear(screen.getByRole('textbox', { name: /search/i }))
    await user.click(await screen.findByRole('option', { name: newInstrument.name }))

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(capturedRequest).toStrictEqual({
      name: newName,
      songId: songId,
      instrumentId: newInstrument.id,
      sectionId: newSongSection.id,
      bandMemberId: settings.defaultBandMember.id
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newName} added!`)).toBeInTheDocument()

    // reset fields
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument() // due to section
    expect(
      screen.getByRole('button', { name: settings.defaultInstrument.name })
    ).toBeInTheDocument()
    expect(screen.getByRole('combobox', { name: 'song-section' })).toHaveValue('')
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
  })

  // Validation

  it('should display error when name was typed and then removed', async () => {
    const user = userEvent.setup()

    const newName = 'New Name'

    reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: {
        songId: '',
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)
    await user.clear(screen.getByRole('textbox', { name: /name/i }))

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()
  })

  it('should not send create request when name is not typed', async () => {
    const user = userEvent.setup()

    let capturedRequest: CreateSongPartRequest
    server.use(
      http.post('/songs/parts', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongPartRequest
        return HttpResponse.json({ message: 'part added!' })
      })
    )

    reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: {
        songId: '',
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()

    expect(capturedRequest).toBeUndefined()
  })

  it('should refresh errors when reopened', async () => {
    const user = userEvent.setup()

    const uut = (opened = true) => <AddNewSongPart opened={opened} onClose={() => {}} />

    const [{ rerender }] = reduxRender(uut(), {
      song: {
        songId: '',
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()

    rerender(uut(false))
    rerender(uut())

    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
  })
})

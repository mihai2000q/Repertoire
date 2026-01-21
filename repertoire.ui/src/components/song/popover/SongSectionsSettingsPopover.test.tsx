import { emptySongSection, emptySongSettings, reduxRender } from '../../../test-utils.tsx'
import SongSectionsSettingsPopover from './SongSectionsSettingsPopover.tsx'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { Instrument, SongSection } from '../../../types/models/Song.ts'
import { BandMember } from '../../../types/models/Artist.ts'
import {
  UpdateAllSongSectionsRequest,
  UpdateSongSettingsRequest
} from '../../../types/requests/SongRequests.ts'

describe('Song Sections Settings Popover', () => {
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

  const handlers = [
    http.get('/songs/instruments', () => {
      return HttpResponse.json(instruments)
    }),
    http.put('/songs/settings', () => {
      return HttpResponse.json({ message: 'it worked' })
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRender(
      <SongSectionsSettingsPopover settings={emptySongSettings} sections={[]} songId={''} />
    )

    expect(screen.getByRole('button', { name: 'settings' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'settings' }))
    expect(await screen.findByRole('dialog')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-instrument' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeDisabled()

    rerender(
      <SongSectionsSettingsPopover
        settings={emptySongSettings}
        sections={[]}
        songId={''}
        bandMembers={bandMembers}
      />
    )

    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-band-member' })).not.toBeDisabled()
  })

  it('should render and display current settings', async () => {
    const user = userEvent.setup()

    const defaultInstrument = instruments[1]
    const defaultBandMember = bandMembers[1]

    reduxRender(
      <SongSectionsSettingsPopover
        settings={{ ...emptySongSettings, defaultInstrument, defaultBandMember }}
        sections={[]}
        songId={''}
        bandMembers={bandMembers}
      />
    )

    await user.click(screen.getByRole('button', { name: 'settings' }))

    expect(screen.getByRole('button', { name: defaultInstrument.name })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: defaultBandMember.name })).toBeInTheDocument()
  })

  describe('band member', () => {
    it('should send update song settings request when changing the band member and display updated popover', async () => {
      const user = userEvent.setup()

      let capturedRequest: UpdateSongSettingsRequest
      server.use(
        http.put('/songs/settings', async (req) => {
          capturedRequest = (await req.request.json()) as UpdateSongSettingsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const songSettings = {
        ...emptySongSettings,
        id: 'some-id'
      }
      const newBandMember = bandMembers[1]

      const sections: SongSection[] = [{ ...emptySongSection, bandMember: bandMembers[0] }]

      reduxRender(
        <SongSectionsSettingsPopover
          settings={songSettings}
          sections={sections}
          songId={''}
          bandMembers={bandMembers}
        />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-band-member' }))
      await user.click(await screen.findByRole('option', { name: newBandMember.name }))

      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          settingsId: songSettings.id,
          defaultBandMemberId: newBandMember.id
        })
      )

      expect(screen.getByText(/update all sections' band members/i)).toBeInTheDocument()
    })

    it('should send update all song sections request when changing the band member and accepting the band members popover', async () => {
      const user = userEvent.setup()

      let capturedRequest: UpdateAllSongSectionsRequest
      server.use(
        http.put('/songs/sections/all', async (req) => {
          capturedRequest = (await req.request.json()) as UpdateAllSongSectionsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const newBandMember = bandMembers[1]
      const songId = 'some-song-id'

      const sections: SongSection[] = [{ ...emptySongSection, bandMember: bandMembers[0] }]

      reduxRender(
        <SongSectionsSettingsPopover
          settings={emptySongSettings}
          sections={sections}
          songId={songId}
          bandMembers={bandMembers}
        />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-band-member' }))
      await user.click(await screen.findByRole('option', { name: newBandMember.name }))

      await user.click(await screen.findByRole('button', { name: 'confirm' }))
      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          songId: songId,
          bandMemberId: newBandMember.id
        })
      )
    })

    it('should not always display the updated band member popover (when there are no sections with distinct members)', async () => {
      const user = userEvent.setup()

      const newBandMember = bandMembers[1]
      const songId = 'some-song-id'

      const sections: SongSection[] = [{ ...emptySongSection, bandMember: newBandMember }]

      reduxRender(
        <SongSectionsSettingsPopover
          settings={emptySongSettings}
          sections={sections}
          songId={songId}
          bandMembers={bandMembers}
        />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-band-member' }))
      await user.click(await screen.findByRole('option', { name: newBandMember.name }))

      expect(screen.queryByText(/update all sections' band members/i)).not.toBeInTheDocument()
    })
  })

  describe('instrument', () => {
    it('should send update song settings request when changing the instrument and display updated popover', async () => {
      const user = userEvent.setup()

      let capturedRequest: UpdateSongSettingsRequest
      server.use(
        http.put('/songs/settings', async (req) => {
          capturedRequest = (await req.request.json()) as UpdateSongSettingsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const songSettings = {
        ...emptySongSettings,
        id: 'some-id'
      }
      const newInstrument = instruments[1]

      const sections: SongSection[] = [{ ...emptySongSection, instrument: instruments[0] }]

      reduxRender(
        <SongSectionsSettingsPopover settings={songSettings} sections={sections} songId={''} />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-instrument' }))
      await user.click(await screen.findByRole('option', { name: newInstrument.name }))

      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          settingsId: songSettings.id,
          defaultInstrumentId: newInstrument.id
        })
      )

      expect(screen.getByText(/update all sections' instruments/i)).toBeInTheDocument()
    })

    it('should send update all song sections request when changing the instrument and accepting the instruments popover', async () => {
      const user = userEvent.setup()

      let capturedRequest: UpdateAllSongSectionsRequest
      server.use(
        http.put('/songs/sections/all', async (req) => {
          capturedRequest = (await req.request.json()) as UpdateAllSongSectionsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const newInstrument = instruments[1]
      const songId = 'some-song-id'

      const sections: SongSection[] = [{ ...emptySongSection, instrument: instruments[0] }]

      reduxRender(
        <SongSectionsSettingsPopover
          settings={emptySongSettings}
          sections={sections}
          songId={songId}
        />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-instrument' }))
      await user.click(await screen.findByRole('option', { name: newInstrument.name }))

      await user.click(await screen.findByRole('button', { name: 'confirm' }))
      await waitFor(() =>
        expect(capturedRequest).toStrictEqual({
          songId: songId,
          instrumentId: newInstrument.id
        })
      )
    })

    it('should not always display the updated instrument popover (when there are no sections with distinct instruments)', async () => {
      const user = userEvent.setup()

      const newInstrument = instruments[1]
      const songId = 'some-song-id'

      const sections: SongSection[] = [{ ...emptySongSection, instrument: newInstrument }]

      reduxRender(
        <SongSectionsSettingsPopover
          settings={emptySongSettings}
          sections={sections}
          songId={songId}
        />
      )

      await user.click(screen.getByRole('button', { name: 'settings' }))
      await user.click(screen.getByRole('button', { name: 'select-instrument' }))
      await user.click(await screen.findByRole('option', { name: newInstrument.name }))

      expect(screen.queryByText(/update all sections' instruments/i)).not.toBeInTheDocument()
    })
  })
})

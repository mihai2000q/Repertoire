import { reduxRender, withToastify } from '../../test-utils.tsx'
import BandMemberCard from './BandMemberCard.tsx'
import { BandMember } from '../../types/models/Artist.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'

describe('Band Member Card', () => {
  const handlers = [
    http.get('/artists/band-members/roles', async () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const bandMember: BandMember = {
    id: '1',
    name: 'Member',
    color: '#123123',
    imageUrl: 'avatar.png',
    roles: []
  }

  it('should render', () => {
    reduxRender(<BandMemberCard bandMember={bandMember} artistId={''} />)

    expect(screen.getByText(bandMember.name)).toBeInTheDocument()
    expect(screen.getByRole('img', { name: bandMember.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: bandMember.name })).toHaveAttribute(
      'src',
      bandMember.imageUrl
    )
  })

  it('should display band member details when clicking on it', async () => {
    const user = userEvent.setup()

    reduxRender(<BandMemberCard bandMember={bandMember} artistId={''} />)

    await user.click(screen.getByRole('img', { name: bandMember.name }))

    expect(await screen.findByRole('dialog')).toBeInTheDocument() // the dialog has no "official" title
    expect(await screen.findByRole('heading', { name: bandMember.name })).toBeInTheDocument()
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRender(<BandMemberCard bandMember={bandMember} artistId={''} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByRole('img', { name: bandMember.name })
    })

    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display edit band member modal when clicking on edit', async () => {
      const user = userEvent.setup()

      reduxRender(<BandMemberCard bandMember={bandMember} artistId={''} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByRole('img', { name: bandMember.name })
      })
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(await screen.findByRole('dialog', { name: /edit band member/i })).toBeInTheDocument()
    })

    it('should display warning modal and then delete when clicking on delete', async () => {
      const user = userEvent.setup()

      const artistId = 'some-id'

      server.use(
        http.delete(`/artists/band-members/${bandMember.id}/from/${artistId}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(withToastify(<BandMemberCard bandMember={bandMember} artistId={artistId} />))

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByRole('img', { name: bandMember.name })
      })
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(screen.getByText(`${bandMember.name} deleted!`)).toBeInTheDocument()
    })
  })
})

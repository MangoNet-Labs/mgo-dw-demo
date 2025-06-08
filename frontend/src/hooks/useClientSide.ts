import { useEffect } from 'react'
import { clientSideRenderringEvent } from '../services/events'

let isInited = false
export default function useClientSide() {
  const isClientSide = typeof window !== 'undefined'
  useEffect(() => {
    if (!isInited && isClientSide) {
      clientSideRenderringEvent()
      isInited = true
    }
  }, [isClientSide])
}

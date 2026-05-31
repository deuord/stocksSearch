import { useRef, useMemo } from 'react'
import { Canvas, useFrame } from '@react-three/fiber'
import * as THREE from 'three'

interface GlobeProps {
  highlightRegion: string
  activeStock: string
}

function WireframeGlobe({ highlightRegion }: { highlightRegion: string }) {
  const meshRef = useRef<THREE.Mesh>(null)
  const linesRef = useRef<THREE.LineSegments>(null)

  const { vertices, indices } = useMemo(() => {
    const radius = 2
    const segments = 48
    const verts: number[] = []
    const idxs: number[] = []

    for (let lat = 0; lat <= segments; lat++) {
      const theta = (lat * Math.PI) / segments
      const sinTheta = Math.sin(theta)
      const cosTheta = Math.cos(theta)

      for (let lon = 0; lon <= segments; lon++) {
        const phi = (lon * 2 * Math.PI) / segments
        const sinPhi = Math.sin(phi)
        const cosPhi = Math.cos(phi)

        const x = cosPhi * sinTheta
        const y = cosTheta
        const z = sinPhi * sinTheta

        verts.push(radius * x, radius * y, radius * z)
      }
    }

    for (let lat = 0; lat < segments; lat++) {
      for (let lon = 0; lon < segments; lon++) {
        const first = lat * (segments + 1) + lon
        const second = first + segments + 1
        idxs.push(first, second)
        idxs.push(first, first + 1)
      }
    }

    const lastRow = segments * (segments + 1)
    for (let lon = 0; lon < segments; lon++) {
      const first = lastRow + lon
      idxs.push(first, first + 1)
      const lastCol = (segments + 1) * (segments + 1) - 1
      for (let lat = 0; lat < segments; lat++) {
        const first = (lat + 1) * (segments + 1) - 1
        idxs.push(first, first + segments + 1)
      }
    }

    return { vertices: new Float32Array(verts), indices: new Uint16Array(idxs) }
  }, [])

  const lineGeometry = useMemo(() => {
    const geo = new THREE.BufferGeometry()
    geo.setAttribute('position', new THREE.BufferAttribute(vertices, 3))
    geo.setIndex(new THREE.BufferAttribute(indices, 1))
    return geo
  }, [vertices, indices])

  const regionColors: Record<string, THREE.Color> = {
    US: new THREE.Color('#4a9eff'),
    HK: new THREE.Color('#ff6b6b'),
    CN: new THREE.Color('#ffd93d'),
    KR: new THREE.Color('#6bcb77'),
    AU: new THREE.Color('#ff922b'),
    JP: new THREE.Color('#cc5de8'),
    SG: new THREE.Color('#20c997'),
  }

  useFrame((_, delta) => {
    if (meshRef.current) {
      meshRef.current.rotation.y += delta * 0.1
    }
    if (linesRef.current) {
      linesRef.current.rotation.y += delta * 0.1
    }
  })

  const glowColor = highlightRegion ? regionColors[highlightRegion] || new THREE.Color('#ffffff') : new THREE.Color('#ffffff')

  return (
    <group ref={meshRef}>
      <lineSegments ref={linesRef} geometry={lineGeometry}>
        <lineBasicMaterial color="#333333" transparent opacity={0.4} />
      </lineSegments>
      <mesh>
        <sphereGeometry args={[2.02, 48, 48]} />
        <meshBasicMaterial color="#000000" transparent opacity={0} />
      </mesh>

      {/* 高亮流光环 */}
      {highlightRegion && (
        <>
          <mesh rotation={[Math.PI / 2, 0, 0]}>
            <torusGeometry args={[2.1, 0.02, 16, 100]} />
            <meshBasicMaterial color={glowColor} transparent opacity={0.8} />
          </mesh>
          <mesh rotation={[Math.PI / 4, 0, Math.PI / 6]}>
            <torusGeometry args={[2.15, 0.01, 16, 100, Math.PI]} />
            <meshBasicMaterial color={glowColor} transparent opacity={0.4} />
          </mesh>
        </>
      )}

      {/* 外层粒子光环 - 流光效果 */}
      <ParticleRing color={glowColor} active={!!highlightRegion} />
    </group>
  )
}

function ParticleRing({ color, active }: { color: THREE.Color; active: boolean }) {
  const pointsRef = useRef<THREE.Points>(null)

  const particles = useMemo(() => {
    const count = 300
    const positions = new Float32Array(count * 3)
    for (let i = 0; i < count; i++) {
      const angle = (i / count) * Math.PI * 2
      const radius = 2.25 + Math.random() * 0.15
      positions[i * 3] = Math.cos(angle) * radius
      positions[i * 3 + 1] = (Math.random() - 0.5) * 0.5
      positions[i * 3 + 2] = Math.sin(angle) * radius
    }
    return positions
  }, [])

  const particleGeo = useMemo(() => {
    const geo = new THREE.BufferGeometry()
    geo.setAttribute('position', new THREE.BufferAttribute(particles, 3))
    return geo
  }, [particles])

  useFrame((_, delta) => {
    if (pointsRef.current && active) {
      pointsRef.current.rotation.y += delta * 0.3
      pointsRef.current.rotation.x += delta * 0.05
    }
  })

  return (
    <points ref={pointsRef} geometry={particleGeo}>
      <pointsMaterial
        color={color}
        size={0.03}
        transparent
        opacity={active ? 0.7 : 0.1}
        blending={THREE.AdditiveBlending}
        depthWrite={false}
      />
    </points>
  )
}

export default function Globe({ highlightRegion }: GlobeProps) {
  return (
    <Canvas
      camera={{ position: [0, 0, 6], fov: 45 }}
      style={{ background: 'transparent' }}
      gl={{ antialias: true, alpha: true }}
    >
      <ambientLight intensity={0.2} />
      <WireframeGlobe highlightRegion={highlightRegion} />
    </Canvas>
  )
}

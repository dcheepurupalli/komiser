import Graph from 'react-graph-vis';
import { useState, useMemo } from 'react';
import { Node } from 'vis';
import { Tooltip as ReactTooltip } from 'react-tooltip';

export type Vertices = {
  id: string;
  name: string;
  service: string;
  color: string;
  data: any;
};

export type Edges = {
  from: string;
  to: string;
  name: string;
};

type InfrastrucuteGraphProps = {
  vertices: Vertices[];
  edges: Edges[];
};

const options = {
  layout: {
    hierarchical: false
  },
  edges: {
    color: '#000000'
  },
  interaction: {
    hover: true,
    tooltipDelay: 100,
    hoverConnectedEdges: true
  },
  nodes: {
    shape: 'dot',
    size: 12,
    scaling: {
      label: true
    },
    shadow: true
  },
  autoResize: true,
  groups: {
    useDefaultGroups: true,
    EC2: { color: { background: 'red' }, borderWidth: 3 },
    'Network Interface': { color: { background: 'green' }, borderWidth: 3 },
    'Security Group': { color: { background: 'pink' }, borderWidth: 3 }
  }
};

function randomColor() {
  const red = Math.floor(Math.random() * 256)
    .toString(16)
    .padStart(2, '0');
  const green = Math.floor(Math.random() * 256)
    .toString(16)
    .padStart(2, '0');
  const blue = Math.floor(Math.random() * 256)
    .toString(16)
    .padStart(2, '0');
  return `#${red}${green}${blue}`;
}

function nodeVertices(
  vertices: Vertices[],
  edges: Edges[]
): { withEdges: Node[]; withoutEdges: Node[] } {
  // Get all the edges from and to into single array.
  const nodes = edges.reduce((acc, edge) => {
    acc.push(edge.from);
    acc.push(edge.to);
    return acc;
  }, [] as string[]);

  const uniqueNodes = [...new Set(nodes)];

  // Seperate vertices in uniqueNode to one array and not int unique nodes to different array.
  const withEdges: Node[] = [];
  const withoutEdges: Node[] = [];
  vertices.forEach(vertex => {
    if (uniqueNodes.includes(vertex.id)) {
      withEdges.push({
        id: vertex.id,
        label: vertex.service,
        group: vertex.service,
        title: vertex.name
      });
    } else {
      withoutEdges.push({
        id: vertex.id,
        label: vertex.service,
        group: vertex.service
      });
    }
  });
  return { withEdges, withoutEdges };
}

function InfrastructureGraph({ vertices, edges }: InfrastrucuteGraphProps) {
  const { withEdges, withoutEdges } = nodeVertices(vertices, edges);

  const [state, setState] = useState<{
    counter: number;
    graph: { nodes: Node[]; edges: Edges[] };
    events: {
      // select: (event: { nodes: Node[]; edges: Edges[] }) => void;
      // click: (event: {}) => void;
      // doubleClick: (event: {
      //   pointer: { canvas: { x: number; y: number } };
      // }) => void;
    };
  }>({
    counter: 10,
    graph: {
      nodes: withEdges,
      edges
    },
    events: {
      // select: ({ nodes, edges }) => {
      // },
      // click: event => {
      // },
      // doubleClick: ({ pointer: { canvas } }) => {
      //   createNode(canvas.x, canvas.y);
      // }
    }
  });

  // const createNode = (x: any, y: any) => {
  //   const color = randomColor();
  //   setState(({ graph: { nodes, edges }, counter, ...rest }) => {
  //     const id = counter + 1;
  //     const from = Math.floor(Math.random() * counter) + 1;
  //     return {
  //       graph: {
  //         nodes: [...nodes, { id, label: `Node ${id}`, color, x, y }],
  //         edges: [...edges, { from, to: id }]
  //       },
  //       counter: id,
  //       ...rest
  //     };
  //   });
  // };

  const { graph } = state;

  const events = {
    // select: (event: any) => {
    //   var { nodes, edges } = event;
    //   if (nodes.length > 0) {
    //     const selectedNode = withEdges.find(node => node.id === nodes[0]);
    //   }
    // },
    // click: function(event:any) {
    //     console.log(event)
    // },
    // showPopup: function(event:any) {
    //     console.log(event)
    // },
    // doubleClick: (event: any) => {
    //   const {
    //     pointer: { canvas }
    //   } = event;
    //   createNode(canvas.x, canvas.y);
    // }
  };
  return (
    <div>
      <Graph
        graph={graph}
        options={options}
        events={events}
        style={{ height: '640px' }}
      />
    </div>
  );
}

export default InfrastructureGraph;

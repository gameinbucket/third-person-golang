package main

import (
	g "./graphics"
)

const (
	bipedScale        = 0.03
	bipedBody         = 0
	bipedHead         = 1
	bipedLeftArm      = 2
	bipedLeftForearm  = 3
	bipedRightArm     = 4
	bipedRightForearm = 5
	bipedLeftLeg      = 6
	bipedLeftKnee     = 7
	bipedRightLeg     = 8
	bipedRightKnee    = 9
	bipedBones        = 10
)

func makeBiped(o *g.Model) {
	o.Bones = make([]g.Bone, bipedBones)

	/* torso */

	o.Bones[bipedBody].Width = 8 * bipedScale
	o.Bones[bipedBody].Height = 16 * bipedScale
	o.Bones[bipedBody].Length = 6 * bipedScale

	/* head */

	o.Bones[bipedHead].Width = 6 * bipedScale
	o.Bones[bipedHead].Height = 6 * bipedScale
	o.Bones[bipedHead].Length = 6 * bipedScale
	o.Bones[bipedHead].BoneOffsetX = 0
	o.Bones[bipedHead].BoneOffsetY = o.Bones[bipedBody].Height + o.Bones[bipedHead].Height
	o.Bones[bipedHead].BoneOffsetZ = 0

	/* arms */

	o.Bones[bipedLeftArm].Width = 3 * bipedScale
	o.Bones[bipedLeftArm].Height = 8 * bipedScale
	o.Bones[bipedLeftArm].Length = 3 * bipedScale
	o.Bones[bipedLeftArm].BoneOffsetX = o.Bones[bipedBody].Width + o.Bones[bipedLeftArm].Width
	o.Bones[bipedLeftArm].BoneOffsetY = o.Bones[bipedBody].Height * 0.8
	o.Bones[bipedLeftArm].BoneOffsetZ = 0
	o.Bones[bipedLeftArm].PlaneOffsetX = 0
	o.Bones[bipedLeftArm].PlaneOffsetY = -(o.Bones[bipedLeftArm].Height)
	o.Bones[bipedLeftArm].PlaneOffsetZ = 0

	o.Bones[bipedLeftForearm].Width = 3 * bipedScale
	o.Bones[bipedLeftForearm].Height = 8 * bipedScale
	o.Bones[bipedLeftForearm].Length = 3 * bipedScale
	o.Bones[bipedLeftForearm].BoneOffsetX = 0
	o.Bones[bipedLeftForearm].BoneOffsetY = -(o.Bones[bipedLeftArm].Height * 2)
	o.Bones[bipedLeftForearm].BoneOffsetZ = 0
	o.Bones[bipedLeftForearm].PlaneOffsetX = 0
	o.Bones[bipedLeftForearm].PlaneOffsetY = -(o.Bones[bipedLeftForearm].Height)
	o.Bones[bipedLeftForearm].PlaneOffsetZ = 0

	o.Bones[bipedRightArm].Width = 3 * bipedScale
	o.Bones[bipedRightArm].Height = 8 * bipedScale
	o.Bones[bipedRightArm].Length = 3 * bipedScale
	o.Bones[bipedRightArm].BoneOffsetX = -(o.Bones[bipedBody].Width + o.Bones[bipedRightArm].Width)
	o.Bones[bipedRightArm].BoneOffsetY = o.Bones[bipedBody].Height * 0.8
	o.Bones[bipedRightArm].BoneOffsetZ = 0
	o.Bones[bipedRightArm].PlaneOffsetX = 0
	o.Bones[bipedRightArm].PlaneOffsetY = -(o.Bones[bipedRightArm].Height)
	o.Bones[bipedRightArm].PlaneOffsetZ = 0

	o.Bones[bipedRightForearm].Width = 3 * bipedScale
	o.Bones[bipedRightForearm].Height = 8 * bipedScale
	o.Bones[bipedRightForearm].Length = 3 * bipedScale
	o.Bones[bipedRightForearm].BoneOffsetX = 0
	o.Bones[bipedRightForearm].BoneOffsetY = -(o.Bones[bipedRightArm].Height * 2)
	o.Bones[bipedRightForearm].BoneOffsetZ = 0
	o.Bones[bipedRightForearm].PlaneOffsetX = 0
	o.Bones[bipedRightForearm].PlaneOffsetY = -(o.Bones[bipedRightForearm].Height)
	o.Bones[bipedRightForearm].PlaneOffsetZ = 0

	/* legs */

	o.Bones[bipedLeftLeg].Width = 3 * bipedScale
	o.Bones[bipedLeftLeg].Height = 8 * bipedScale
	o.Bones[bipedLeftLeg].Length = 3 * bipedScale
	o.Bones[bipedLeftLeg].BoneOffsetX = o.Bones[bipedBody].Width - o.Bones[bipedLeftLeg].Width
	o.Bones[bipedLeftLeg].BoneOffsetY = -(o.Bones[bipedBody].Height)
	o.Bones[bipedLeftLeg].BoneOffsetZ = 0
	o.Bones[bipedLeftLeg].PlaneOffsetX = 0
	o.Bones[bipedLeftLeg].PlaneOffsetY = -(o.Bones[bipedLeftLeg].Height)
	o.Bones[bipedLeftLeg].PlaneOffsetZ = 0

	o.Bones[bipedLeftKnee].Width = 3 * bipedScale
	o.Bones[bipedLeftKnee].Height = 8 * bipedScale
	o.Bones[bipedLeftKnee].Length = 3 * bipedScale
	o.Bones[bipedLeftKnee].BoneOffsetX = 0
	o.Bones[bipedLeftKnee].BoneOffsetY = -(o.Bones[bipedLeftLeg].Height * 2)
	o.Bones[bipedLeftKnee].BoneOffsetZ = 0
	o.Bones[bipedLeftKnee].PlaneOffsetX = 0
	o.Bones[bipedLeftKnee].PlaneOffsetY = -(o.Bones[bipedLeftKnee].Height)
	o.Bones[bipedLeftKnee].PlaneOffsetZ = 0

	o.Bones[bipedRightLeg].Width = 3 * bipedScale
	o.Bones[bipedRightLeg].Height = 8 * bipedScale
	o.Bones[bipedRightLeg].Length = 3 * bipedScale
	o.Bones[bipedRightLeg].BoneOffsetX = -(o.Bones[bipedBody].Width - o.Bones[bipedRightLeg].Width)
	o.Bones[bipedRightLeg].BoneOffsetY = -(o.Bones[bipedBody].Height)
	o.Bones[bipedRightLeg].BoneOffsetZ = 0
	o.Bones[bipedRightLeg].PlaneOffsetX = 0
	o.Bones[bipedRightLeg].PlaneOffsetY = -(o.Bones[bipedRightLeg].Height)
	o.Bones[bipedRightLeg].PlaneOffsetZ = 0

	o.Bones[bipedRightKnee].Width = 3 * bipedScale
	o.Bones[bipedRightKnee].Height = 8 * bipedScale
	o.Bones[bipedRightKnee].Length = 3 * bipedScale
	o.Bones[bipedRightKnee].BoneOffsetX = 0
	o.Bones[bipedRightKnee].BoneOffsetY = -(o.Bones[bipedRightLeg].Height * 2)
	o.Bones[bipedRightKnee].BoneOffsetZ = 0
	o.Bones[bipedRightKnee].PlaneOffsetX = 0
	o.Bones[bipedRightKnee].PlaneOffsetY = -(o.Bones[bipedRightKnee].Height)
	o.Bones[bipedRightKnee].PlaneOffsetZ = 0

	/* leafs */

	o.Bones[bipedBody].Leafs = make([]*g.Bone, 5)
	o.Bones[bipedBody].Leafs[0] = &o.Bones[bipedHead]
	o.Bones[bipedBody].Leafs[1] = &o.Bones[bipedLeftArm]
	o.Bones[bipedBody].Leafs[2] = &o.Bones[bipedRightArm]
	o.Bones[bipedBody].Leafs[3] = &o.Bones[bipedLeftLeg]
	o.Bones[bipedBody].Leafs[4] = &o.Bones[bipedRightLeg]

	o.Bones[bipedLeftArm].Leafs = make([]*g.Bone, 1)
	o.Bones[bipedLeftArm].Leafs[0] = &o.Bones[bipedLeftForearm]

	o.Bones[bipedRightArm].Leafs = make([]*g.Bone, 1)
	o.Bones[bipedRightArm].Leafs[0] = &o.Bones[bipedRightForearm]

	o.Bones[bipedLeftLeg].Leafs = make([]*g.Bone, 1)
	o.Bones[bipedLeftLeg].Leafs[0] = &o.Bones[bipedLeftKnee]

	o.Bones[bipedRightLeg].Leafs = make([]*g.Bone, 1)
	o.Bones[bipedRightLeg].Leafs[0] = &o.Bones[bipedRightKnee]

	/* roots */

	o.Bones[bipedBody].RecursiveInit(nil)
}

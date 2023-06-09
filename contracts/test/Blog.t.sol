// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "../src/Blog.sol";

contract BlogTest is Test {
    Blog public blog;

    function setUp() public {
        blog = new Blog();
    }

    function testCreatePost() public {
        string memory cid = "cid_string";
        string memory title = "TITLE";
        blog.createPost(title, cid);

        uint256 postId = blog.getCurrentPostId();
        Blog.Post memory post = blog.getPost(postId - 1);
    
        assertEq(post.title, title, "incorrect title");
        assertEq(post.cid, cid, "incorrect CID");
        assertEq(post.author, address(this), "incorrect author");
    }

}

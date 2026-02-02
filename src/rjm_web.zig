const std = @import("std");

pub fn downloadFile(allocator: std.mem.Allocator, url: []const u8, output_path: []const u8) !void {
    var client = std.http.Client{ .allocator = allocator };
    defer client.deinit();

    const uri = try std.Uri.parse(url);

    var redirect_buffer: [8 * 1024]u8 = undefined;

    // Create output file
    const file = try std.fs.cwd().createFile(output_path, .{});
    defer file.close();

    var file_writer_buffer: [8 * 1024]u8 = undefined;
    var file_writer = file.writer(&file_writer_buffer);

    std.debug.print("Downloading {s}...\n", .{url});

    const result = try client.fetch(.{
        .location = .{ .uri = uri },
        .method = .GET,
        .redirect_buffer = &redirect_buffer,
        .response_writer = &file_writer.interface,
        .headers = .{
            .user_agent = .{ .override = "zig-download-example" },
        },
    });

    try file_writer.interface.flush();

    if (result.status != .ok) {
        std.debug.print("Download failed with status: {}\n", .{result.status});
        return error.HttpError;
    }

    std.debug.print("Download complete: {s}\n", .{output_path});
}

// const ProgressWriter = struct {
//     child_writer: std.fs.File.Writer,
//     downloaded_bytes: u64,
//     last_mb: u64,

//     pub fn write(self: *ProgressWriter, bytes: []const u8) !usize {
//         const written = try self.child_writer.write(bytes);
//         self.downloaded_bytes += written;

//         // Update progress every MB
//         const current_mb = self.downloaded_bytes / (1024 * 1024);
//         if (current_mb != self.last_mb) {
//             self.last_mb = current_mb;
//             self.printProgress();
//         }

//         return written;
//     }

//     fn printProgress(self: *ProgressWriter) void {
//         const mb_downloaded = @as(f64, @floatFromInt(self.downloaded_bytes)) / (1024.0 * 1024.0);
//         std.debug.print("\rDownloaded: {d:.2} MB", .{mb_downloaded});
//     }
// };

pub fn downloadFileWithProgress(allocator: std.mem.Allocator, url: []const u8, output_path: []const u8) !void {
    var client = std.http.Client{ .allocator = allocator };
    defer client.deinit();
    const uri = try std.Uri.parse(url);
    var redirect_buffer: [8 * 1024]u8 = undefined;

    // Create output file
    const file = try std.fs.cwd().createFile(output_path, .{});
    defer file.close();
    var file_writer_buffer: [8 * 1024]u8 = undefined;
    var file_writer = file.writer(&file_writer_buffer);

    // Progress tracking
    var downloaded_bytes: u64 = 0;
    var last_mb: u64 = 0;

    const Context = struct {
        file_writer: *std.fs.File.Writer,
        downloaded_bytes: *u64,
        last_mb: *u64,
    };

    var ctx = Context{
        .file_writer = &file_writer,
        .downloaded_bytes = &downloaded_bytes,
        .last_mb = &last_mb,
    };

    var progress_writer = std.io.AnyWriter{
        .context = &ctx,
        .writeFn = struct {
            fn writeFn(context: *const anyopaque, bytes: []const u8) anyerror!usize {
                const c: *Context = @constCast(@ptrCast(@alignCast(context)));
                const written = try c.file_writer.interface.write(bytes);
                c.downloaded_bytes.* += written;

                const current_mb = c.downloaded_bytes.* / (1024 * 1024);
                if (current_mb != c.last_mb.*) {
                    c.last_mb.* = current_mb;
                    const mb = @as(f64, @floatFromInt(c.downloaded_bytes.*)) / (1024.0 * 1024.0);
                    std.debug.print("\rDownloaded: {d:.2} MB", .{mb});
                }

                return written;
            }
        }.writeFn,
    };

    std.debug.print("Downloading {s}...\n", .{url});

    const result = try client.fetch(.{
        .location = .{ .uri = uri },
        .method = .GET,
        .redirect_buffer = &redirect_buffer,
        .response_writer = &progress_writer,
        .headers = .{
            .user_agent = .{ .override = "zig-download-example" },
        },
    });

    try file_writer.interface.flush();

    if (result.status != .ok) {
        std.debug.print("\nDownload failed with status: {}\n", .{result.status});
        return error.HttpError;
    }

    std.debug.print("\nDownload complete: {s}\n", .{output_path});
}

import argparse
import json
import math
from PIL import Image
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.progress import Progress

console = Console()

# --- Core Logic ---

def image_to_vector(image_path):
    try:
        with Image.open(image_path) as img:
            img = img.convert("RGB")
            pixels = list(img.getdata())
            return [p for pixel in pixels for p in pixel]
    except FileNotFoundError:
        console.print(f"[bold red]Error: Image file not found at {image_path}[/bold red]")
        return None
    except Exception as e:
        console.print(f"[bold red]An error occurred: {e}[/bold red]")
        return None

def save_vector(vector_data, path, file_format):
    try:
        with open(path, 'w') as f:
            if file_format == 'json':
                json.dump(vector_data, f, indent=4)
            else:
                f.write('\n'.join(map(str, vector_data)))
        return True
    except Exception as e:
        console.print(f"[bold red]Error saving vector file: {e}[/bold red]")
        return False

def read_vector(path, file_format):
    try:
        with open(path, 'r') as f:
            if file_format == 'json':
                return json.load(f)
            else:
                return [int(line.strip()) for line in f if line.strip()]
    except FileNotFoundError:
        console.print(f"[bold red]Error: Vector file not found at {path}[/bold red]")
        return None
    except Exception as e:
        console.print(f"[bold red]Error reading vector file: {e}[/bold red]")
        return None

def vector_to_image(vector_data, width, height, output_path):
    expected_len = width * height * 3
    if len(vector_data) != expected_len:
        console.print(f"[bold red]Vector size mismatch! Expected {expected_len}, got {len(vector_data)}.[/bold red]")
        return False
    try:
        pixels = list(zip(*[iter(vector_data)] * 3))
        with Image.new('RGB', (width, height)) as img:
            img.putdata(pixels)
            img.save(output_path)
        return True
    except Exception as e:
        console.print(f"[bold red]Error creating image: {e}[/bold red]")
        return False

# --- TUI Handlers ---

def handle_to_vector(args):
    console.print(Panel(f"Converting [bold cyan]{args.input}[/bold cyan] to a numerical vector...", title="Image to Vector", expand=False))
    vector = image_to_vector(args.input)
    if vector:
        with Progress(console=console) as progress:
            task = progress.add_task("[green]Saving vector...", total=1)
            if save_vector(vector, args.output, args.format):
                progress.update(task, advance=1)
                console.print(f"[bold green]Success![/bold green] Vector saved to [cyan]{args.output}[/cyan] ({args.format.upper()}).")
                console.print(f"  - Vector Size: {len(vector)} values")

def handle_to_image(args):
    console.print(Panel(f"Converting [bold cyan]{args.input}[/bold cyan] back to an image...", title="Vector to Image", expand=False))
    vector = read_vector(args.input, args.format)
    if vector:
        with Progress(console=console) as progress:
            task = progress.add_task("[green]Reconstructing image...", total=1)
            if vector_to_image(vector, args.width, args.height, args.output):
                progress.update(task, advance=1)
                console.print(f"[bold green]Success![/bold green] Image saved to [cyan]{args.output}[/cyan].")

def handle_check(args):
    console.print(Panel("Checking Dimensions", title="Dimension Check", expand=False))
    if args.image:
        try:
            with Image.open(args.image) as img:
                w, h = img.size
                table = Table(title=f"Image Details: {args.image}")
                table.add_column("Property", style="cyan")
                table.add_column("Value", style="magenta")
                table.add_row("Dimensions", f"{w} x {h}")
                table.add_row("Total Pixels", f"{w*h:,}")
                table.add_row("RGB Vector Size", f"{w*h*3:,}")
                console.print(table)
        except FileNotFoundError:
            console.print(f"[bold red]Error: Image file not found at {args.image}[/bold red]")

    if args.vector:
        vector = read_vector(args.vector, args.format)
        if vector:
            size = len(vector)
            table = Table(title=f"Vector Details: {args.vector}")
            table.add_column("Property", style="cyan")
            table.add_column("Value", style="magenta")
            table.add_row("Format", args.format.upper())
            table.add_row("Vector Size", f"{size:,} values")
            console.print(table)
            if size % 3 == 0:
                pixels = size // 3
                dim_table = Table(title=f"Possible RGB Dimensions for {pixels:,} Pixels")
                dim_table.add_column("Width", style="cyan")
                dim_table.add_column("Height", style="magenta")
                for i in range(1, int(math.sqrt(pixels)) + 1):
                    if pixels % i == 0:
                        w, h = i, pixels // i
                        dim_table.add_row(str(w), str(h))
                        if w != h: dim_table.add_row(str(h), str(w))
                console.print(dim_table)
            else:
                console.print("[yellow]Vector size is not divisible by 3; cannot determine RGB dimensions.[/yellow]")

# --- Main Argparse Setup ---

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="A versatile tool for image and vector conversions.",
        formatter_class=argparse.RawTextHelpFormatter
    )
    subparsers = parser.add_subparsers(dest='command', required=True, help='Available commands')

    # 'convert' command
    p_convert = subparsers.add_parser('convert', help='Convert between images and vectors.', epilog=" ")
    conv_sub = p_convert.add_subparsers(dest='direction', required=True)
    
    # to-vector
    p_to_vec = conv_sub.add_parser('to-vector', help='Image to Vector', 
        epilog='Example:\n  python img-vec-tool.py convert to-vector i.jpg vector.txt --format txt')
    p_to_vec.add_argument('input', help='Input image path')
    p_to_vec.add_argument('output', help='Output vector file path')
    p_to_vec.add_argument('--format', choices=['txt', 'json'], default='txt', help='Output format')
    p_to_vec.set_defaults(func=handle_to_vector)

    # to-image
    p_to_img = conv_sub.add_parser('to-image', help='Vector to Image', 
        epilog='Example:\n  python img-vec-tool.py convert to-image vector.txt new_image.png 3072 4096')
    p_to_img.add_argument('input', help='Input vector file path')
    p_to_img.add_argument('output', help='Output image path')
    p_to_img.add_argument('width', type=int, help='Image width')
    p_to_img.add_argument('height', type=int, help='Image height')
    p_to_img.add_argument('--format', choices=['txt', 'json'], default='txt', help='Input format')
    p_to_img.set_defaults(func=handle_to_image)

    # 'check' command
    p_check = subparsers.add_parser('check', help='Check dimensions of images or vectors.', 
        epilog=('Examples:\n'
                '  python img-vec-tool.py check --image i.jpg\n'
                '  python img-vec-tool.py check --vector vector.txt\n'
                '  python img-vec-tool.py check --image i.jpg --vector vector.txt'))
    p_check.add_argument('-i', '--image', help='Path to the image file')
    p_check.add_argument('-v', '--vector', help='Path to the vector file')
    p_check.add_argument('--format', choices=['txt', 'json'], default='txt', help='Vector file format')
    p_check.set_defaults(func=handle_check)

    args = parser.parse_args()
    if hasattr(args, 'func'):
        args.func(args)
    else:
        parser.print_help()
